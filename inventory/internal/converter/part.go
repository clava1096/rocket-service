package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/clava1096/rocket-service/inventory/internal/model"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

func ValueFromProto(v *inventoryv1.Value) model.Value {
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

func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func PartToProto(p model.Part) *inventoryv1.Part {
	// Metadata
	metadata := make(map[string]*inventoryv1.Value, len(p.Metadata))
	for k, v := range p.Metadata {
		metadata[k] = ValueToProto(v)
	}

	// Category
	var category inventoryv1.Category
	switch p.Category {
	case model.CategoryEngine:
		category = inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		category = inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPortholes:
		category = inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		category = inventoryv1.Category_CATEGORY_WING
	default:
		category = inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
	return &inventoryv1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: int64(p.StockQuantity),
		Category:      category,
		Dimensions: &inventoryv1.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  metadata,
		CreatedAt: timeToTimestampByValue(p.CreatedAt),
		UpdatedAt: timeToTimestamp(p.UpdatedAt),
	}
}

func ValueToProto(v model.Value) *inventoryv1.Value {
	switch v.Kind {
	case model.ValueKindString:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_StringValue{StringValue: v.StringValue}}
	case model.ValueKindInt64:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_Int64Value{Int64Value: v.IntegerValue}}
	case model.ValueKindFloat64:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_DoubleValue{DoubleValue: v.DoubleValue}}
	case model.ValueKindBool:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_BoolValue{BoolValue: v.BooleanValue}}
	default:
		return nil
	}
}

func timeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func timeToTimestampByValue(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func FilterFromProto(filter *inventoryv1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	// Конвертируем категории
	categories := make([]model.Category, len(filter.Categories))
	for i, cat := range filter.Categories {
		switch cat {
		case inventoryv1.Category_CATEGORY_ENGINE:
			categories[i] = model.CategoryEngine
		case inventoryv1.Category_CATEGORY_FUEL:
			categories[i] = model.CategoryFuel
		case inventoryv1.Category_CATEGORY_PORTHOLE:
			categories[i] = model.CategoryPortholes
		case inventoryv1.Category_CATEGORY_WING:
			categories[i] = model.CategoryWing
		default:
			// Игнорируем CATEGORY_UNSPECIFIED или неизвестные
		}
	}

	return model.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
