package user

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/iam/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.User, error) {
	user, err := s.userRepository.Get(ctx, uuid)

	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return model.User{}, model.ErrUserAlreadyExists
		}
		return model.User{}, err
	}

	return user, nil
}
