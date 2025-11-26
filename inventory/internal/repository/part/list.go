package part

import (
	"context"
	"github.com/clava1096/rocket-service/inventory/internal/model"
	"github.com/clava1096/rocket-service/inventory/internal/repository/converter"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part
	for _, part := range r.inventory {
		parts = append(parts, converter.PartFromRepoModel(part))
	}

	if len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0 {
		return parts, nil
	}

	if len(filter.Uuids) > 0 {
		parts = r.filterByUUIDs(parts, filter.Uuids)
	}
	if len(filter.Names) > 0 {
		parts = r.filterByNames(parts, filter.Names)
	}
	if len(filter.Categories) > 0 {
		parts = r.filterByCategories(parts, filter.Categories)
	}
	if len(filter.ManufacturerCountries) > 0 {
		parts = r.filterByManufacturerCountries(parts, filter.ManufacturerCountries)
	}
	if len(filter.Tags) > 0 {
		parts = r.filterByTags(parts, filter.Tags)
	}

	return parts, nil
}

func (r *repository) filterByUUIDs(parts []model.Part, uuids []string) []model.Part {
	uuidSet := make(map[string]struct{}, len(uuids))
	for _, id := range uuids {
		uuidSet[id] = struct{}{}
	}

	var filtered []model.Part
	for _, part := range parts {
		if _, found := uuidSet[part.Uuid]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func (r *repository) filterByNames(parts []model.Part, names []string) []model.Part {
	nameSet := make(map[string]struct{}, len(names))
	for _, name := range names {
		nameSet[name] = struct{}{}
	}

	var filtered []model.Part
	for _, part := range parts {
		if _, found := nameSet[part.Name]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func (r *repository) filterByCategories(parts []model.Part, categories []model.Category) []model.Part {
	categorySet := make(map[model.Category]struct{}, len(categories))
	for _, cat := range categories {
		categorySet[cat] = struct{}{}
	}

	var filtered []model.Part
	for _, part := range parts {
		if _, found := categorySet[part.Category]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func (r *repository) filterByManufacturerCountries(parts []model.Part, countries []string) []model.Part {
	countrySet := make(map[string]struct{}, len(countries))
	for _, country := range countries {
		countrySet[country] = struct{}{}
	}

	var filtered []model.Part
	for _, part := range parts {
		if _, found := countrySet[part.Manufacturer.Country]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func (r *repository) filterByTags(parts []model.Part, tags []string) []model.Part {
	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}

	var filtered []model.Part
	for _, part := range parts {
		for _, tag := range part.Tags {
			if _, found := tagSet[tag]; found {
				filtered = append(filtered, part)
				break // достаточно одного совпадения
			}
		}
	}
	return filtered
}
