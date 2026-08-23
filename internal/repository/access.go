package repository

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type Access interface {
	CountUsers(context.Context) (int64, error)
	CreateUser(context.Context, domain.User, string) error
	FindUserByEmail(context.Context, string) (domain.User, string, error)
	GetUser(context.Context, string) (domain.User, error)
	ListUsers(context.Context) ([]domain.User, error)
	UpdateUserRole(context.Context, string, domain.Role) error
	SetUserActive(context.Context, string, bool) error
	MarkLogin(context.Context, string, time.Time) error
	CreateSession(context.Context, domain.Session) error
	FindSession(context.Context, string, time.Time) (domain.Session, domain.User, error)
	TouchSession(context.Context, string, time.Time) error
	DeleteSession(context.Context, string) error
	DeleteExpiredSessions(context.Context, time.Time) (int64, error)
}
