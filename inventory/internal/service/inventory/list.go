package inventory

import (
	"context"

	"github.com/clava1096/rocket-service/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.List(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}
	return parts, nil
}
