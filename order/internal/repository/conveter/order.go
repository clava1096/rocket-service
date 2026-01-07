package converter

import (
	"github.com/google/uuid"

	"github.com/clava1096/rocket-service/order/internal/model"
	repoModel "github.com/clava1096/rocket-service/order/internal/repository/model"
)

// OrderToRepoModel конвертирует доменную модель Order в репозиторную.
// Глубоко копирует все вложенные структуры и коллекции для изоляции слоёв.
func OrderToRepoModel(order model.Order) repoModel.Order {
	// Копируем опциональное поле TransactionUUID
	var txUUID *string
	if order.TransactionUUID != nil {
		uuidCopy := *order.TransactionUUID
		txUUID = &uuidCopy
	}

	// Копируем опциональное поле PaymentMethod
	var payMethod *repoModel.PaymentMethod
	if order.PaymentMethod != nil {
		methodCopy := repoModel.PaymentMethod(*order.PaymentMethod)
		payMethod = &methodCopy
	}

	orderUuid, _ := uuid.Parse(order.UUID)
	userUuid, _ := uuid.Parse(order.UserUUID)
	return repoModel.Order{
		UUID:            orderUuid,
		UserUUID:        userUuid,
		PartUUIDs:       copyUUIDSlice(order.PartUUIDs), // глубокая копия среза
		TotalPrice:      order.TotalPrice,
		Status:          repoModel.OrderStatus(order.Status),
		TransactionUUID: txUUID,
		PaymentMethod:   payMethod,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

// OrderFromRepoModel конвертирует репозиторную модель Order в доменную.
func OrderFromRepoModel(order repoModel.Order) model.Order {
	var txUUID *string
	if order.TransactionUUID != nil {
		uuidCopy := *order.TransactionUUID
		txUUID = &uuidCopy
	}

	var payMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		methodCopy := model.PaymentMethod(*order.PaymentMethod)
		payMethod = &methodCopy
	}

	return model.Order{
		UUID:            order.UUID.String(),
		UserUUID:        order.UserUUID.String(),
		PartUUIDs:       copyStringSlice(order.PartUUIDs),
		TotalPrice:      order.TotalPrice,
		Status:          model.OrderStatus(order.Status),
		TransactionUUID: txUUID,
		PaymentMethod:   payMethod,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

// copyUUIDSlice создаёт глубокую копию среза UUID.
// Сохраняет семантику nil: если входной срез nil — возвращает nil.
func copyUUIDSlice(src []string) []uuid.UUID {
	if src == nil {
		return nil
	}

	dst := make([]uuid.UUID, len(src))
	for i := range src {
		dst[i], _ = uuid.Parse(src[i])
	}
	return dst
}

// copyStringSlice создаёт глубокую копию среза строк.
// Сохраняет семантику nil: если входной срез nil — возвращает nil.
func copyStringSlice(src []uuid.UUID) []string {
	if src == nil {
		return nil
	}

	dst := make([]string, len(src))
	for i := range src {
		dst[i] = src[i].String()
	}
	return dst
}
