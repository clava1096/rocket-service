package v1

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/iam/internal/model"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Get(ctx context.Context, req *iamV1.GetUserRequest) (*iamV1.GetUserResponse, error) {
	user, err := a.iamService.Get(ctx, req.UserUuid)
	if err != nil {
		switch {

		case errors.Is(err, model.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")

		default:
			return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
		}
	}

	notifMethods := make([]*iamV1.NotificationMethod, 0, len(user.NotificationMethods))
	for _, nm := range user.NotificationMethods {
		notifMethods = append(notifMethods, &iamV1.NotificationMethod{
			ProviderName: nm.ProviderName,
			Target:       nm.Target,
		})
	}

	return &iamV1.GetUserResponse{
		UserUuid:            user.ID.String(),
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notifMethods,
	}, nil
}
