package decoder

import (
	"github.com/clava1096/rocket-service/order/internal/model"
	eventsv1 "github.com/clava1096/rocket-service/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder { return &orderPaidDecoder{} }

func (d *orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, err
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
