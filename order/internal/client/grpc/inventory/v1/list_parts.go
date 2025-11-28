package v1

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/client/converter"
	"github.com/clava1096/rocket-service/order/internal/model"

	generatedInventoryV1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

// ListParts реализует метод из интерфейса InventoryClient.
func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {

	pbFilter := converter.PartsFilterToProto(filter)

	resp, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: pbFilter,
	})
	if err != nil {
		return nil, err
	}

	return converter.PartListToModel(resp.Parts), nil
}
