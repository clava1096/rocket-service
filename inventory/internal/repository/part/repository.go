package part

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// repo "github.com/clava1096/rocket-service/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/clava1096/rocket-service/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database, coll string) *repository {
	collection := db.Collection(coll)

	r := &repository{
		collection: collection,
	}

	if err := r.createIndexes(context.Background()); err != nil {
		log.Printf("WARNING: failed to create indexes: %v", err)
	}

	r.initSampleData()
	return r
}

func (r *repository) createIndexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "manufacturer.country", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}
