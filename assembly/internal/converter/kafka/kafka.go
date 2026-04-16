package kafka

import "github.com/clava1096/rocket-service/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode([]byte) (model.OrderPaidEvent, error)
}
