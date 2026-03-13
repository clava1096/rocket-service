package app

import (
	"context"
	"fmt"

	"github.com/clava1096/rocket-service/inventory/internal/config"
	inventoryRepository "github.com/clava1096/rocket-service/inventory/internal/repository/part"
	inventoryService "github.com/clava1096/rocket-service/inventory/internal/service/inventory"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryAPI "github.com/clava1096/rocket-service/inventory/internal/api/inventory/v1"
	"github.com/clava1096/rocket-service/inventory/internal/repository"
	"github.com/clava1096/rocket-service/inventory/internal/service"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1         inventoryv1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryAPI(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryV1 == nil {
		d.inventoryV1 = inventoryAPI.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryV1
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx), config.AppConfig().MongoDb.CollectionName())
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().MongoDb.URI()))

		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}
		err = client.Ping(ctx, readpref.Primary())

		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().MongoDb.DatabaseName())
	}

	return d.mongoDBHandle
}
