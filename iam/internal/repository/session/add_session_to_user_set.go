package session

import (
	"context"
	"fmt"

	"github.com/clava1096/rocket-service/iam/internal/config"
	"github.com/google/uuid"
)

func (r *repository) AddSessionToUserSet(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	setKey := config.AppConfig().Session.MakeUserSessionsKey(userID.String())

	err := r.cacheClient.SAdd(ctx, setKey, sessionID.String())
	if err != nil {
		return fmt.Errorf("failed to add session to user set: %w", err)
	}

	return nil
}
