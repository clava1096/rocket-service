package session

import (
	"context"
	"time"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/cache"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, ttl time.Duration) error
	Get(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	AddSessionToUserSet(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
}

type repository struct {
	cacheClient cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cacheClient: cache,
	}
}
