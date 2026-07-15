package user

import (
	"context"
	"fmt"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, req model.RegisterRequest) (uuid.UUID, error) {
	exists, err := s.userRepository.ExistsByLoginOrEmail(ctx, req.Login, req.Email)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return uuid.Nil, model.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := model.User{
		ID:                  uuid.New(),
		Login:               req.Login,
		Email:               req.Email,
		PasswordHash:        string(passwordHash),
		NotificationMethods: req.NotificationMethods,
	}

	_, err = s.userRepository.Create(ctx, newUser)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser.ID, nil
}
