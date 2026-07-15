package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/clava1096/rocket-service/iam/internal/repository/converter"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")

func (r *repository) Get(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	value, err := r.cacheClient.Get(ctx, r.sessionConfig.KeyPrefix()+sessionID.String())

	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from cache: %w", err)
	}

	if len(value) == 0 {
		return nil, ErrSessionNotFound
	}

	userID, err := converter.ValueToUserId(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert value to user id: %w", err)
	}

	return &model.Session{
		ID:     sessionID,
		UserID: userID,
	}, nil
}
