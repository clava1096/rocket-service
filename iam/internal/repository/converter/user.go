package converter

import (
	"github.com/clava1096/rocket-service/iam/internal/model"
	repoModel "github.com/clava1096/rocket-service/iam/internal/repository/model"
)

func UserFromRepoModel(user repoModel.User) model.User {
	return model.User{
		ID:                  user.UUID,
		Login:               user.Login,
		Email:               user.Email,
		PasswordHash:        user.PasswordHash,
		NotificationMethods: NotificationMethodsRepoModel(user.NotificationMethods),
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func NotificationMethodsRepoModel(notificationMethods []repoModel.NotificationMethod) []model.NotificationMethod {
	if len(notificationMethods) == 0 {
		return []model.NotificationMethod{}
	}

	dst := make([]model.NotificationMethod, len(notificationMethods))

	for i := range notificationMethods {
		dst[i] = notificationMethodRepoModel(notificationMethods[i])
	}
	return dst
}

func notificationMethodRepoModel(notificationMethod repoModel.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}

func NotificationMethodsDomainModel(notificationMethods []model.NotificationMethod) []repoModel.NotificationMethod {
	if len(notificationMethods) == 0 {
		return []repoModel.NotificationMethod{}
	}

	dst := make([]repoModel.NotificationMethod, len(notificationMethods))

	for i := range notificationMethods {
		dst[i] = notificationMethodDomainModel(notificationMethods[i])
	}
	return dst
}

func notificationMethodDomainModel(notificationMethod model.NotificationMethod) repoModel.NotificationMethod {
	return repoModel.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}
