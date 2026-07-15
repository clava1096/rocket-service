package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (*model.User, error) {
	sessionID, err := uuid.Parse(sessionUUID)
	if err != nil {
		return nil, model.ErrInvalidSessionFormat
	}

	session, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, model.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	user, err := s.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
