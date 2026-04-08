package service

import (
	"context"

	"github.com/clava1096/rocket-service/notification/internal/model"
)

type OrderPaidConsumer interface {
	//OrderPaidHandler(ctx context.Context, msg kafka.Message) error
	RunConsumer(ctx context.Context) error
}

type OrderAssemblyConsumer interface {
	//OrderAssemblyHandler(ctx context.Context, msg kafka.Message) error
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error
}
