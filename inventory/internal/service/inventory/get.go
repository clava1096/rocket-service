package inventory

import (
	"context"
	"github.com/clava1096/rocket-service/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, uuid)

	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
