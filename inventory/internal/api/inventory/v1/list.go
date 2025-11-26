package v1

import (
	"context"
	"github.com/clava1096/rocket-service/inventory/internal/converter"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.inventoryService.List(ctx, converter.FilterFromProto(req.Filter))
	if err != nil {
		return nil, err
	}

	protoParts := make([]*inventoryv1.Part, len(parts))

	for i, part := range parts {
		protoParts[i] = converter.PartToProto(part)
	}

	return &inventoryv1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
