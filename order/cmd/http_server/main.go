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

	orderAPI "github.com/clava1096/rocket-service/order/internal/api/order/v1"
	inventoryv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/inventory/v1"
	paymentv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/payment/v1"
	"github.com/clava1096/rocket-service/order/internal/config"
	m "github.com/clava1096/rocket-service/order/internal/middleware"
	"github.com/clava1096/rocket-service/order/internal/migrator"
	orderRepository "github.com/clava1096/rocket-service/order/internal/repository/order"
	orderService "github.com/clava1096/rocket-service/order/internal/service/order"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()
	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatal("error while loading config:", err)
	}

	runApp(ctx, cfg)
}

func runApp(ctx context.Context, cfg *config.Config) {
	pool := initDatabase(cfg, ctx)
	defer func() {
		pool.Close()
	}()

	initMigration(cfg)
	inventoryConn, rawInventoryClient := initInventory(cfg)
	inventoryClient := inventoryv1Client.NewClient(rawInventoryClient)
	defer inventoryConn.Close()

	paymentConn, rawPaymentClient := initPayment(cfg)
	paymentClient := paymentv1Client.NewClient(rawPaymentClient)
	defer paymentConn.Close()

	repo := orderRepository.NewRepository(pool)
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
		Addr:              net.JoinHostPort("localhost", cfg.Server.HttpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Println("Starting server on " + cfg.Server.HttpPort)
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

func initDatabase(cfg *config.Config, ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.Postgres.GetPostgresUri())

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	return pool
}

func initMigration(cfg *config.Config) {
	pgxConfig, err := pgx.ParseConfig(cfg.Postgres.GetPostgresUri())
	if err != nil {
		log.Fatalf("Error parsing database URI: %v", err)
	}

	log.Println("Running migrations...")
	sqlDb := stdlib.OpenDB(*pgxConfig)
	defer sqlDb.Close()

	migratorRunner := migrator.NewMigrator(sqlDb, cfg.Postgres.MigrationsDir)
	if err = migratorRunner.Up(); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("done")
}

func initInventory(cfg *config.Config) (*grpc.ClientConn, inventoryv1.InventoryServiceClient) {
	connInventory, err := grpc.NewClient(
		cfg.Server.InventoryGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}

	rawInventoryClient := inventoryv1.NewInventoryServiceClient(connInventory)

	return connInventory, rawInventoryClient
}

func initPayment(cfg *config.Config) (*grpc.ClientConn, paymentv1.PaymentServiceClient) {
	connPayment, err := grpc.NewClient(
		cfg.Server.PaymentGrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	rawPaymentClient := paymentv1.NewPaymentServiceClient(connPayment)
	return connPayment, rawPaymentClient
}
