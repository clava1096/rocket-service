package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clava1096/rocket-service/inventory/internal/converter"
	"github.com/clava1096/rocket-service/inventory/internal/model"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, err := a.inventoryService.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part: %s not found", req.Uuid)
		}
		return nil, err
	}

	protoPart := converter.PartToProto(part)

	return &inventoryv1.GetPartResponse{
		Part: protoPart,
	}, nil
}
