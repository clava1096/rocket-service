package v1

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/iam/internal/model"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Whoami(ctx context.Context, req *iamV1.WhoamiRequest) (*iamV1.WhoamiResponse, error) {

	usr, err := a.sessionService.Whoami(ctx, req.SessionUuid)

	if err != nil {
		switch {

		case errors.Is(err, model.ErrInvalidSessionFormat):
			return nil, status.Error(codes.InvalidArgument, "invalid session uuid format")

		case errors.Is(err, model.ErrSessionNotFound):
			return nil, status.Error(codes.Unauthenticated, "invalid or expired session")

		case errors.Is(err, model.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")

		default:
			return nil, status.Errorf(codes.Internal, "failed to validate session: %v", err)
		}
	}

	return &iamV1.WhoamiResponse{
		UserUuid: usr.ID.String(),
		Login:    usr.Login,
		Email:    usr.Email,
	}, nil
}
