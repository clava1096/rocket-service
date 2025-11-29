package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/clava1096/rocket-service/order/internal/api/order/v1"
	inventoryv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/inventory/v1"
	paymentv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/payment/v1"
	m "github.com/clava1096/rocket-service/order/internal/middleware"
	orderRepository "github.com/clava1096/rocket-service/order/internal/repository/order"
	orderService "github.com/clava1096/rocket-service/order/internal/service/order"
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

func main() {
	connInventory, err := grpc.NewClient(
		inventoryGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}
	defer connInventory.Close()

	rawInventoryClient := inventoryv1.NewInventoryServiceClient(connInventory)
	inventoryClient := inventoryv1Client.NewClient(rawInventoryClient)

	connPayment, err := grpc.NewClient(
		paymentGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	defer connPayment.Close()

	rawPaymentClient := paymentv1.NewPaymentServiceClient(connPayment)
	paymentClient := paymentv1Client.NewClient(rawPaymentClient)
	repo := orderRepository.NewRepository()
	service := orderService.NewService(repo, inventoryClient, paymentClient)
	api := orderAPI.NewAPI(service)

	storageServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("Error creating OpenApi server: %v", err)
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
		log.Println("Starting server on " + httpPort)
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server was stopped.")
}
