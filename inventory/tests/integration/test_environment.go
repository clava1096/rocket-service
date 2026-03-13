package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/clava1096/rocket-service/inventory/internal/model"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (env *TestEnvironment) InsertTestDetail(ctx context.Context) (string, error) {
	uuid := gofakeit.UUID()
	now := time.Now()
	weight := gofakeit.Float64()
	metadata := map[string]model.Value{
		"color": {
			Kind:        model.ValueKindString,
			StringValue: "red",
		},
		"weight_kg": {
			Kind:        model.ValueKindFloat64,
			DoubleValue: weight,
		},
	}

	part := bson.M{
		"_id":            uuid,
		"name":           gofakeit.Name(),
		"description":    gofakeit.Phrase(),
		"price":          gofakeit.Float64(),
		"stock_quantity": gofakeit.Float64(),
		"category":       model.CategoryPortholes,
		"dimensions": bson.M{
			"length": gofakeit.Float64(),
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"weight": weight,
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags":       []string{gofakeit.Word(), gofakeit.Word()},
		"metadata":   metadata,
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")

	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (env *TestEnvironment) InsertTestDetailsWithData(ctx context.Context, part *inventoryv1.Part) (string, error) {
	now := time.Now()
	partToMongo := bson.M{
		"_id":            part.Uuid,
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       part.Category,
		"dimensions": bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		},
		"manufacturer": bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		},
		"tags":       part.Tags,
		"metadata":   part.Metadata,
		"created_at": primitive.NewDateTimeFromTime(part.CreatedAt.AsTime()),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")

	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, partToMongo)
	if err != nil {
		return "", err
	}

	return part.Uuid, nil
}

func (env *TestEnvironment) GetTestPartInfo() *inventoryv1.Part {
	now := time.Now()
	weight := gofakeit.Float64Range(0.1, 50)

	return &inventoryv1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Phrase(),
		Price:         gofakeit.Price(10, 10000),
		StockQuantity: int64(gofakeit.Int64()),
		Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryv1.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Word(), gofakeit.Word()},
		Metadata: map[string]*inventoryv1.Value{
			"color": {
				Kind: &inventoryv1.Value_StringValue{
					StringValue: "red",
				},
			},
			"weight_kg": {
				Kind: &inventoryv1.Value_DoubleValue{
					DoubleValue: weight,
				},
			},
		},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
}

// GetUpdatedPartInfo — возвращает обновленную информацию о детали (для тестов Update)
func (env *TestEnvironment) GetUpdatedPartInfo() *inventoryv1.Part {
	// Используем wrapperspb для опциональных полей в update-запросах
	return &inventoryv1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Phrase(),
		Price:         gofakeit.Price(10, 10000),
		StockQuantity: gofakeit.Int64(),
		Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryv1.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 100),
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Word(), gofakeit.Word()},
	}
}

// ClearInventoryCollection — удаляет все записи из коллекции inventory (parts)
func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
