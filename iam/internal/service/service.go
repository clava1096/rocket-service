package service

import (
	"context"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"
)

type IAMService interface {
	Get(ctx context.Context, uuid string) (model.User, error)

	Register(ctx context.Context, req model.RegisterRequest) (uuid.UUID, error)
}

type SessionService interface {
	Login(ctx context.Context, login string, password string) (uuid.UUID, error)

	Whoami(ctx context.Context, sessionUUID string) (*model.User, error)
}
