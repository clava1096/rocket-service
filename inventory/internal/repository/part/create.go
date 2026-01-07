package part

import (
	"context"

	"github.com/clava1096/rocket-service/inventory/internal/model"
	"github.com/clava1096/rocket-service/inventory/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, part model.Part) (model.Part, error) {
	repoPart := converter.PartToRepoModel(part)

	_, err := r.collection.InsertOne(ctx, repoPart)
	if err != nil {
		return model.Part{}, err // model.ErrWhileCreate
	}

	return part, nil
}
