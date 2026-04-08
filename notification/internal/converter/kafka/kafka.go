package kafka

import "github.com/clava1096/rocket-service/notification/internal/model"

type OrderPaidDecoder interface {
	Decode([]byte) (model.OrderPaidEvent, error)
}

type OrderAssemblyDecoder interface {
	Decode([]byte) (model.ShipAssembledEvent, error)
}
