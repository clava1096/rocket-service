package repository

import (
	"context"

	"github.com/clava1096/rocket-service/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)

	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)

	Create(ctx context.Context, part model.Part) (model.Part, error)
}
