package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	m "github.com/clava1096/rocket-service/order/internal/middleware"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	inventoryGrpcPort = ":50051"
	paymentGrpcPort   = ":50052"
)

type OrderStorage struct {
	mu       sync.Mutex
	orderDto map[string]*orderV1.OrderDto
}

type OrderHandler struct {
	store           *OrderStorage
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orderDto: make(map[string]*orderV1.OrderDto),
	}
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryv1.InventoryServiceClient, paymentClient paymentv1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		store:           storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (s *OrderStorage) getOrder(uuid string) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()
	orderDto, ok := s.orderDto[uuid]
	if !ok {
		return nil
	}
	return orderDto
}

func (s *OrderStorage) delOrder(order *orderV1.OrderDto) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.Status = orderV1.OrderStatusCANCELLED

	return order
}

func (s *OrderStorage) newOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orderDto[order.OrderUUID.String()] = order
	return nil
}

func (s *OrderStorage) updateOrder(order *orderV1.OrderDto) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orderDto[order.OrderUUID.String()] = order
	return order
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderUUID := params.OrderUUID
	order := h.store.getOrder(orderUUID)
	if order == nil {
		return &orderV1.NotFoundError{Message: "Order '" + orderUUID + "' not found"}, nil
	}
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Message: "Order '" + orderUUID + "' was paid"}, nil
	}
	order = h.store.delOrder(order)
	if order == nil {
		return &orderV1.NotFoundError{Message: "Order '" + orderUUID + "' not found"}, nil
	}
	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) CreateOrderRequest(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRequestRes, error) {
	listPart, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartsUUID,
		},
	})
	if err != nil {
		log.Printf("Inventory service call failed: %v", err)
		return &orderV1.InternalServerError{
			Message: "Something went wrong",
		}, nil
	}

	found := make(map[string]struct{}, len(listPart.Parts))
	for _, p := range listPart.Parts {
		found[p.GetUuid()] = struct{}{}
	}
	for _, id := range req.PartsUUID {
		if _, ok := found[id]; !ok {
			return &orderV1.BadRequestError{
				Message: fmt.Sprintf("part not found: %s", id),
			}, nil
		}
	}

	var totalPrice float64
	for _, part := range listPart.Parts {
		totalPrice += part.Price
	}

	partUUIDs := make([]uuid.UUID, len(req.PartsUUID))
	for i, s := range req.PartsUUID {
		u, err := uuid.Parse(s)
		if err != nil {
			return &orderV1.BadRequestError{
				Message: fmt.Sprintf("invalid part UUID: %s", s),
			}, nil
		}
		partUUIDs[i] = u
	}
	order := &orderV1.OrderDto{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.MustParse(req.UUID),
		PartUuids:  partUUIDs,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
		TotalPrice: totalPrice,
	}
	if err = h.store.newOrder(order); err != nil {
		log.Printf("Failed to save order %s: %v", order.GetOrderUUID(), err)
		return &orderV1.InternalServerError{
			Message: "Failed to create order",
		}, nil
	}
	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.GetOrderUUID(),
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetInfoOrderByUUID(_ context.Context, params orderV1.GetInfoOrderByUUIDParams) (orderV1.GetInfoOrderByUUIDRes, error) {
	orderUUID := params.OrderUUID
	order := h.store.getOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order :'" + orderUUID + "' was not found!",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) OrderPayment(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPaymentParams) (orderV1.OrderPaymentRes, error) {
	orderUUID := params.OrderUUID
	order := h.store.getOrder(orderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order :'" + orderUUID + "' was not found!",
		}, nil
	}

	paymentMethod, err := convertPaymentMethod(req.PaymentMethod)
	if err != nil {
		return &orderV1.BadRequestError{
			Message: err.Error(),
		}, nil
	}

	payResp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.GetOrderUUID().String(),
		UserUuid:      order.GetUserUUID().String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Message: "Something went wrong",
		}, nil
	}

	txnUUID, err := uuid.Parse(payResp.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Message: "Invalid transaction UUID from payment service",
		}, nil
	}

	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptNilUUID(txnUUID)
	order.PaymentMethod = orderV1.NewOptPaymentMethod(req.PaymentMethod)

	h.store.updateOrder(order)

	return &orderV1.PayOrderResponse{TransactionUUID: uuid.MustParse(payResp.TransactionUuid)}, nil
}

func convertPaymentMethod(openAPI orderV1.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch openAPI {
	case orderV1.PaymentMethodUNKNOWN:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, nil
	case orderV1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderV1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderV1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("unsupported payment method: %s", openAPI)
	}
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("Application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	storage := NewOrderStorage()

	// Подключение к InventoryService
	connInventory, err := grpc.NewClient(
		"localhost"+inventoryGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to inventory service: %v", err)
		return err
	}
	defer func() {
		if cerr := connInventory.Close(); cerr != nil {
			log.Printf("Failed to close inventory connection: %v", cerr)
		}
	}()

	// Подключение к PaymentService
	connPayment, err := grpc.NewClient(
		"localhost"+paymentGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to payment service: %v", err)
		return err
	}
	defer func() {
		if cerr := connPayment.Close(); cerr != nil {
			log.Printf("Failed to close payment connection: %v", cerr)
		}
	}()

	inventoryClient := inventoryv1.NewInventoryServiceClient(connInventory)
	paymentClient := paymentv1.NewPaymentServiceClient(connPayment)
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	storageServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Error creating OpenAPI server: %v", err)
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(m.RequestLogger)
	r.Mount("/", storageServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Starting server on %s", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
		return err
	}
	log.Println("Server was stopped.")
	return nil
}
