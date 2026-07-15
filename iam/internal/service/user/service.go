package user

import (
	"github.com/clava1096/rocket-service/iam/internal/repository"
	def "github.com/clava1096/rocket-service/iam/internal/service"
)

var _ def.IAMService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
