package decoder

import (
	"github.com/clava1096/rocket-service/notification/internal/model"
	eventsv1 "github.com/clava1096/rocket-service/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type shipAssembledDecoder struct{}

func NewShipAssembledDecoder() *shipAssembledDecoder {
	return &shipAssembledDecoder{}
}

func (d *shipAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, err
	}
	return model.ShipAssembledEvent{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
