package inventory

import "github.com/clava1096/rocket-service/inventory/internal/repository"

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{inventoryRepository}
}
