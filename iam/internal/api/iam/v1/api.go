package v1

import (
	"github.com/clava1096/rocket-service/iam/internal/service"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
)

type api struct {
	iamV1.UnimplementedUserServiceServer

	iamService     service.IAMService
	sessionService service.SessionService
}

func NewAPI(iamService service.IAMService, sessionService service.SessionService) *api {
	return &api{
		iamService:     iamService,
		sessionService: sessionService,
	}
}
