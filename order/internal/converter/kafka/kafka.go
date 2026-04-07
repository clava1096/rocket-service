package kafka

import "github.com/clava1096/rocket-service/order/internal/model"

type OrderPaidDecoder interface {
	Decode([]byte) (model.OrderPaidEvent, error)
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
