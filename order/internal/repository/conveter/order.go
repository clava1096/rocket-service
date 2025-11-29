package converter

import (
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

	return repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       copyStringSlice(order.PartUUIDs), // глубокая копия среза
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
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       copyStringSlice(order.PartUUIDs),
		TotalPrice:      order.TotalPrice,
		Status:          model.OrderStatus(order.Status),
		TransactionUUID: txUUID,
		PaymentMethod:   payMethod,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

// copyStringSlice создаёт глубокую копию среза строк.
// Сохраняет семантику nil: если входной срез nil — возвращает nil.
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
