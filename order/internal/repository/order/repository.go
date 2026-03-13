package model

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/clava1096/rocket-service/order/internal/repository"
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
