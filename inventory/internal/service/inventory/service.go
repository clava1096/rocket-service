package inventory

import (
	"github.com/clava1096/rocket-service/inventory/internal/repository"
	def "github.com/clava1096/rocket-service/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{inventoryRepository}
}
