package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/clava1096/rocket-service/platform/pkg/grpc"
	authV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

type IAMClient = authV1.UserServiceClient

type AuthMiddleware struct {
	iamClient IAMClient
}

func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		whoamiRes, err := m.iamClient.Whoami(r.Context(), &authV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.UserUuid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
