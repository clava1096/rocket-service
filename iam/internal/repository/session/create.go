package session

import (
	"context"
	"fmt"
	"time"

	"github.com/clava1096/rocket-service/iam/internal/config"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, ttl time.Duration) error {
	key := config.AppConfig().Session.MakeKey(sessionID.String())
	value := userID.String()

	err := r.cacheClient.SetWithTTL(ctx, key, value, ttl)
	if err != nil {
		return fmt.Errorf("failed to store sessions in redis: %w", err)
	}

	return nil
}
