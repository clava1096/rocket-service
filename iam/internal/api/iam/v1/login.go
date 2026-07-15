package v1

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/iam/internal/model"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Login(ctx context.Context, loginOrder *iamV1.LoginRequest) (*iamV1.LoginResponse, error) {
	sessionID, err := a.sessionService.Login(ctx, loginOrder.Login, loginOrder.Password)
	if err != nil {

		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}

		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}

		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}

	return &iamV1.LoginResponse{
		SessionUuid: sessionID.String(),
	}, nil
}
