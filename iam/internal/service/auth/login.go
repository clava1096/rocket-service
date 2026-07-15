package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/clava1096/rocket-service/iam/internal/config"
	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, login string, password string) (uuid.UUID, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if errors.Is(err, model.ErrUserNotFound) {
		// 🔐 Не уточняем, что именно не так (логин или пароль) — безопасность!
		return uuid.Nil, model.ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}

	sessionID := uuid.New()

	sessionTTL := config.AppConfig().Session.TTL()

	err = s.sessionRepository.Create(ctx, sessionID, user.ID, sessionTTL)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create session in redis: %w", err)
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.ID, sessionID)
	if err != nil {
		//logger.Logger(ctx, "failed to add session to user set", zap.Error(err))
	}

	return sessionID, nil
}
