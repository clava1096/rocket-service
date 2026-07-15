package v1

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/iam/internal/model"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(ctx context.Context, req *iamV1.RegisterRequest) (*iamV1.RegisterResponse, error) {

	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	notifMethods := convertNotificationMethods(req.GetNotificationMethods())

	usr := model.RegisterRequest{
		Login:               req.Login,
		Password:            req.Password,
		Email:               req.Email,
		NotificationMethods: notifMethods,
	}

	user, err := a.iamService.Register(ctx, usr)

	if err != nil {
		return nil, mapRegisterError(err)
	}

	return &iamV1.RegisterResponse{
		UserUuid: user.String(),
	}, nil

}

func mapRegisterError(err error) error {
	switch {

	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user with this login or email already exists")

	default:
		return status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

}

func validateRegisterRequest(req *iamV1.RegisterRequest) error {

	if req.GetLogin() == "" {

		return status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {

		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetEmail() == "" {

		return status.Error(codes.InvalidArgument, "email is required")
	}
	return nil
}

func convertNotificationMethods(protoMethods []*iamV1.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(protoMethods))

	for _, protoMethod := range protoMethods {
		if protoMethod == nil {
			continue
		}

		result = append(result, model.NotificationMethod{
			ProviderName: protoMethod.GetProviderName(),
			Target:       protoMethod.GetTarget(),
		})
	}

	return result
}
