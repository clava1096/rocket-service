package service

import (
	"context"

	"github.com/clava1096/rocket-service/assembly/internal/model"
)

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
