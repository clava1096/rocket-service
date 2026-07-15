package repository

import (
	"context"
	"time"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Get(ctx context.Context, uuid string) (model.User, error)
	ExistsByLoginOrEmail(ctx context.Context, login, email string) (bool, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	Create(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, ttl time.Duration) error
	AddSessionToUserSet(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
}
