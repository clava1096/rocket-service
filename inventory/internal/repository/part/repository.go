package part

import (
	"sync"

	def "github.com/clava1096/rocket-service/inventory/internal/repository"
	repo "github.com/clava1096/rocket-service/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu        sync.RWMutex
	inventory map[string]repo.Part
}

func NewRepository() *repository {
	r := &repository{
		inventory: make(map[string]repo.Part),
	}
	r.initSampleData()
	return r
}
