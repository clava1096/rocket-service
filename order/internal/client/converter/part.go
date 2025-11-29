package converter

import (
	"github.com/clava1096/rocket-service/order/internal/model"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

// PartsFilterToProto конвертирует доменный фильтр в Protobuf.
func PartsFilterToProto(filter model.PartsFilter) *inventoryv1.PartsFilter {
	return &inventoryv1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categoriesToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// categoriesToProto конвертирует категории.
func categoriesToProto(categories []model.Category) []inventoryv1.Category {
	result := make([]inventoryv1.Category, len(categories))
	for i, cat := range categories {
		result[i] = categoryToProto(cat)
	}
	return result
}

// categoryToProto конвертирует одну категорию.
func categoryToProto(cat model.Category) inventoryv1.Category {
	switch cat {
	case model.CategoryEngine:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPortholes:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

// PartListToModel конвертирует список Protobuf-частей в доменные модели.
func PartListToModel(parts []*inventoryv1.Part) []model.Part {
	result := make([]model.Part, len(parts))
	for i, p := range parts {
		result[i] = PartToModel(p)
	}
	return result
}

// PartToModel конвертирует одну Protobuf-часть в доменную модель.
func PartToModel(p *inventoryv1.Part) model.Part {
	if p == nil {
		return model.Part{}
	}

	metadata := make(map[string]model.Value, len(p.Metadata))
	for k, v := range p.Metadata {
		metadata[k] = ValueToModel(v)
	}

	return model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: float64(p.StockQuantity),
		Category:      categoryFromProto(p.Category),
		Dimensions: model.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  metadata,
		CreatedAt: p.CreatedAt.AsTime(),
		UpdatedAt: p.UpdatedAt.AsTime(),
	}
}

// categoryFromProto конвертирует Protobuf-категорию в доменную.
func categoryFromProto(cat inventoryv1.Category) model.Category {
	switch cat {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPortholes
	case inventoryv1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

// ValueToModel конвертирует Protobuf-Value в доменную модель.
func ValueToModel(v *inventoryv1.Value) model.Value {
	if v == nil {
		return model.Value{}
	}
	switch val := v.Kind.(type) {
	case *inventoryv1.Value_StringValue:
		return model.Value{Kind: model.ValueKindString, StringValue: val.StringValue}
	case *inventoryv1.Value_Int64Value:
		return model.Value{Kind: model.ValueKindInt64, IntegerValue: val.Int64Value}
	case *inventoryv1.Value_DoubleValue:
		return model.Value{Kind: model.ValueKindFloat64, DoubleValue: val.DoubleValue}
	case *inventoryv1.Value_BoolValue:
		return model.Value{Kind: model.ValueKindBool, BooleanValue: val.BoolValue}
	default:
		return model.Value{}
	}
}
