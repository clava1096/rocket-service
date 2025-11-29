package inventory

import (
	"github.com/clava1096/rocket-service/inventory/internal/repository"
	def "github.com/clava1096/rocket-service/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.PartRepository
}

func NewService(inventoryRepository repository.PartRepository) *service {
	return &service{inventoryRepository}
}
