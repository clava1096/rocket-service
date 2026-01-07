package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/clava1096/rocket-service/inventory/internal/api/inventory/v1"
	"github.com/clava1096/rocket-service/inventory/internal/config"
	logger "github.com/clava1096/rocket-service/inventory/internal/middleware"
	inventoryRepository "github.com/clava1096/rocket-service/inventory/internal/repository/part"
	inventoryService "github.com/clava1096/rocket-service/inventory/internal/service/inventory"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

func main() {
	ctx := context.Background()
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("error while init config:", err)
	}
	runApp(cfg, ctx)
}

func runApp(cfg *config.Config, ctx context.Context) {
	lis, err := net.Listen("tcp", cfg.Server.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.LoggerInterceptor()))
	client, db, err := initMongoDB(cfg, ctx)
	if err != nil {
		log.Fatalf("failed to init mongodb: %v", err)
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect: %v", cerr)
		}
	}()

	repo := inventoryRepository.NewRepository(db, cfg.MongoDb.Collection) // todo тут должно быть подклбючение
	service := inventoryService.NewService(repo)
	api := inventoryAPI.NewAPI(service)
	inventoryv1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("Starting gRPC server at %s", cfg.Server.GrpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}

func initMongoDB(cfg *config.Config, ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	dbUri := cfg.MongoDb.GetMongodbUri()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.MongoDb.Database)
	return client, db, nil
}
