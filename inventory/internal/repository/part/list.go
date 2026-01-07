package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/clava1096/rocket-service/inventory/internal/model"
	"github.com/clava1096/rocket-service/inventory/internal/repository/converter"
	repoModel "github.com/clava1096/rocket-service/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	var repoParts []repoModel.Part

	// log.Printf("filter: %v", filter) //todo приходит пустой фильтр

	partsFilter := partsFilters(filter)

	cursor, err := r.collection.Find(ctx, partsFilter)

	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, len(repoParts))

	for i, repoPart := range repoParts {
		parts[i] = converter.PartFromRepoModel(repoPart)
	}

	return parts, nil
}

func partsFilters(filter model.PartsFilter) bson.M {
	query := bson.M{}

	if len(filter.Uuids) > 0 {
		query["_id"] = bson.M{"$in": filter.Uuids}
	}

	if len(filter.Categories) > 0 {
		cats := make([]string, len(filter.Categories))
		for i, c := range filter.Categories {
			cats[i] = string(c)
		}
		query["category"] = bson.M{"$in": cats}
	}

	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	return query
}
