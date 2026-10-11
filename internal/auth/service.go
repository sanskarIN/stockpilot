package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const (
	MinPasswordBytes = 12
	MaxPasswordBytes = 72
)

type Service struct {
	repo     repository.Access
	ttl      time.Duration
	tokenKey []byte
	now      func() time.Time
}

type LoginResult struct {
	User      domain.User
	Token     string
	ExpiresAt time.Time
}

type Principal struct {
	User      domain.User
	SessionID string
	ExpiresAt time.Time
}

func New(repo repository.Access, ttl time.Duration, sessionSecret string) *Service {
	return &Service{
		repo: repo, ttl: ttl, tokenKey: []byte(sessionSecret),
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) BootstrapAdmin(ctx context.Context, email, displayName, password string) (domain.User, error) {
	count, err := s.repo.CountUsers(ctx)
	if err != nil {
		return domain.User{}, err
	}
	if count != 0 {
		return domain.User{}, fmt.Errorf("%w: initial administrator already exists", domain.ErrConflict)
	}
	return s.CreateUser(ctx, email, displayName, domain.RoleAdmin, password)
}

func (s *Service) CreateUser(ctx context.Context, email, displayName string, role domain.Role, password string) (domain.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return domain.User{}, err
	}
	id, err := idgen.New("usr")
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:          id,
		Email:       strings.ToLower(strings.TrimSpace(email)),
		DisplayName: strings.TrimSpace(displayName),
		Role:        role,
		Active:      true,
	}
	if err := user.Validate(); err != nil {
		return domain.User{}, err
	}
	if err := s.repo.CreateUser(ctx, user, string(passwordHash)); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, passwordHash, err := s.repo.FindUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil {
		// Perform comparable password work even when the account is absent to reduce
		// obvious account-enumeration timing differences.
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$12$w8DfB/.bn1j1eQpYLRkpV.4ZdKVvvZSLtZfJHjDptG4NjdZeZErnO"), []byte(password))
		return LoginResult{}, domain.ErrForbidden
	}
	if !user.Active {
		return LoginResult{}, domain.ErrForbidden
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return LoginResult{}, domain.ErrForbidden
	}

	rawToken, err := newSessionToken()
	if err != nil {
		return LoginResult{}, err
	}
	tokenHash := s.hashToken(rawToken)
	now := s.now()
	sessionID, err := idgen.New("ses")
	if err != nil {
		return LoginResult{}, err
	}
	session := domain.Session{
		ID: sessionID, UserID: user.ID, TokenHash: tokenHash,
		CreatedAt: now, LastSeen: now, ExpiresAt: now.Add(s.ttl),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return LoginResult{}, err
	}
	if err := s.repo.MarkLogin(ctx, user.ID, now); err != nil {
		_ = s.repo.DeleteSession(ctx, tokenHash)
		return LoginResult{}, err
	}
	user.LastLoginAt = &now
	return LoginResult{User: user, Token: rawToken, ExpiresAt: session.ExpiresAt}, nil
}

func (s *Service) Resolve(ctx context.Context, rawToken string) (Principal, error) {
	if rawToken == "" {
		return Principal{}, domain.ErrForbidden
	}
	now := s.now()
	session, user, err := s.repo.FindSession(ctx, s.hashToken(rawToken), now)
	if err != nil {
		return Principal{}, domain.ErrForbidden
	}
	if now.Sub(session.LastSeen) >= 5*time.Minute {
		_ = s.repo.TouchSession(ctx, session.ID, now)
	}
	return Principal{User: user, SessionID: session.ID, ExpiresAt: session.ExpiresAt}, nil
}

func (s *Service) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, s.hashToken(rawToken))
}

func (s *Service) PurgeExpired(ctx context.Context) (int64, error) {
	return s.repo.DeleteExpiredSessions(ctx, s.now())
}

func hashPassword(password string) ([]byte, error) {
	if len(password) < MinPasswordBytes || len(password) > MaxPasswordBytes {
		return nil, fmt.Errorf("%w: password must contain %d-%d bytes", domain.ErrInvalid, MinPasswordBytes, MaxPasswordBytes)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	return hash, nil
}

func newSessionToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func (s *Service) hashToken(token string) string {
	mac := hmac.New(sha256.New, s.tokenKey)
	_, _ = mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
