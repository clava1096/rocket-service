package model

import (
	def "github.com/clava1096/rocket-service/order/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}
