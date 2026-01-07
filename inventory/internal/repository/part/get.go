package part

import (
	"context"

	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/clava1096/rocket-service/inventory/internal/model"
	repoConverter "github.com/clava1096/rocket-service/inventory/internal/repository/converter"
	repoModel "github.com/clava1096/rocket-service/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	var part repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrNotFound
		}
	}

	return repoConverter.PartFromRepoModel(part), nil
}
