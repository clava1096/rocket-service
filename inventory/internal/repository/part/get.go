package part

import (
	"context"

	"github.com/clava1096/rocket-service/inventory/internal/model"
	repoConverter "github.com/clava1096/rocket-service/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repoPart, ok := r.inventory[uuid]
	if !ok {
		return model.Part{}, model.ErrNotFound
	}

	return repoConverter.PartFromRepoModel(repoPart), nil
}
