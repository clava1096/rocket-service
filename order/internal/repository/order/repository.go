package model

import (
	"sync"

	def "github.com/clava1096/rocket-service/order/internal/repository"
	repo "github.com/clava1096/rocket-service/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]*repo.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*repo.Order),
	}
}
