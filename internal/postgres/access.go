package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM users`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateUser(ctx context.Context, user domain.User, passwordHash string) error {
	if err := user.Validate(); err != nil {
		return err
	}
	if len(passwordHash) < 20 {
		return fmt.Errorf("%w: password hash is invalid", domain.ErrInvalid)
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO users (id, email, display_name, password_hash, role, active)
		VALUES ($1, lower($2), $3, $4, $5, $6)`,
		user.ID, strings.TrimSpace(user.Email), strings.TrimSpace(user.DisplayName), passwordHash, user.Role, user.Active)
	return mapError(err)
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (domain.User, string, error) {
	var user domain.User
	var passwordHash string
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, display_name, role, active, created_at, updated_at, last_login_at, password_hash
		FROM users WHERE lower(email)=lower($1)`, strings.TrimSpace(email)).Scan(
		&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &passwordHash)
	if err != nil {
		return domain.User{}, "", mapError(err)
	}
	return user, passwordHash, nil
}

func (s *Store) GetUser(ctx context.Context, id string) (domain.User, error) {
	return scanUser(s.pool.QueryRow(ctx, `
		SELECT id, email, display_name, role, active, created_at, updated_at, last_login_at
		FROM users WHERE id=$1`, id))
}

func (s *Store) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, display_name, role, active, created_at, updated_at, last_login_at
		FROM users ORDER BY display_name, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.User, 0)
	for rows.Next() {
		item, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateUserRole(ctx context.Context, id string, role domain.Role) error {
	if !role.Valid() {
		return fmt.Errorf("%w: invalid role", domain.ErrInvalid)
	}
	command, err := s.pool.Exec(ctx, `UPDATE users SET role=$2, updated_at=now() WHERE id=$1`, id, role)
	if err != nil {
		return mapError(err)
	}
	if command.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) SetUserActive(ctx context.Context, id string, active bool) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	command, err := tx.Exec(ctx, `UPDATE users SET active=$2, updated_at=now() WHERE id=$1`, id, active)
	if err != nil {
		return mapError(err)
	}
	if command.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	if !active {
		if _, err := tx.Exec(ctx, `DELETE FROM sessions WHERE user_id=$1`, id); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) MarkLogin(ctx context.Context, id string, at time.Time) error {
	command, err := s.pool.Exec(ctx, `UPDATE users SET last_login_at=$2, updated_at=now() WHERE id=$1`, id, at)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) CreateSession(ctx context.Context, session domain.Session) error {
	if err := session.Validate(time.Now().UTC()); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO sessions (id, user_id, token_hash, expires_at, created_at, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		session.ID, session.UserID, session.TokenHash, session.ExpiresAt, session.CreatedAt, session.LastSeen)
	return mapError(err)
}

func (s *Store) FindSession(ctx context.Context, tokenHash string, now time.Time) (domain.Session, domain.User, error) {
	var session domain.Session
	var user domain.User
	err := s.pool.QueryRow(ctx, `
		SELECT s.id, s.user_id, s.token_hash, s.expires_at, s.created_at, s.last_seen_at,
		       u.id, u.email, u.display_name, u.role, u.active, u.created_at, u.updated_at, u.last_login_at
		FROM sessions s
		JOIN users u ON u.id=s.user_id
		WHERE s.token_hash=$1 AND s.expires_at>$2 AND u.active=true`, tokenHash, now).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt, &session.CreatedAt, &session.LastSeen,
		&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
	if err != nil {
		return domain.Session{}, domain.User{}, mapError(err)
	}
	return session, user, nil
}

func (s *Store) TouchSession(ctx context.Context, id string, at time.Time) error {
	_, err := s.pool.Exec(ctx, `UPDATE sessions SET last_seen_at=$2 WHERE id=$1`, id, at)
	return err
}

func (s *Store) DeleteSession(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash=$1`, tokenHash)
	return err
}

func (s *Store) DeleteExpiredSessions(ctx context.Context, now time.Time) (int64, error) {
	command, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at <= $1`, now)
	if err != nil {
		return 0, err
	}
	return command.RowsAffected(), nil
}

func scanUser(row scanner) (domain.User, error) {
	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt); err != nil {
		return domain.User{}, mapError(err)
	}
	return user, nil
}

var _ repository.Access = (*Store)(nil)
