package v1

import "github.com/clava1096/rocket-service/order/internal/service"

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{orderService}
}
