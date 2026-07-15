package auth

import (
	"github.com/clava1096/rocket-service/iam/internal/repository"
	def "github.com/clava1096/rocket-service/iam/internal/service"
)

var _ def.SessionService = (*service)(nil)

type service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewService(sessionRepository repository.SessionRepository, userRepository repository.UserRepository) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}
