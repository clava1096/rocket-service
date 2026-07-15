package user

import (
	"context"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Get(ctx context.Context, uuid string) (model.User, error)
	ExistsByLoginOrEmail(ctx context.Context, login, email string) (bool, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}
