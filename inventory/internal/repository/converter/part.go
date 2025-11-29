// Package converter обеспечивает преобразование между доменными моделями (internal/model)
// и моделями уровня репозитория (internal/repository/model).
// Это необходимо для строгой изоляции слоёв в архитектуре.
package converter

import (
	"github.com/clava1096/rocket-service/inventory/internal/model"
	repoModel "github.com/clava1096/rocket-service/inventory/internal/repository/model"
)

// ┌───────────────────────────────────────────────────────┐
// │        Конвертация: Репозиторий → Домен               │
// └───────────────────────────────────────────────────────┘

// PartFromRepoModel конвертирует репозиторную модель Part в доменную.
// Глубоко копирует все вложенные структуры и коллекции для изоляции слоёв.
func PartFromRepoModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryFromRepoModel(part.Category),
		Dimensions:    DimensionsFromRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerFromRepoModel(part.Manufacturer),
		Tags:          copyStringSlice(part.Tags),
		Metadata:      MetadataFromRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// ┌───────────────────────────────────────────────────────┐
// │        Конвертация: Домен → Репозиторий               │
// └───────────────────────────────────────────────────────┘

// PartToRepoModel конвертирует доменную модель Part в репозиторную.
// Глубоко копирует все вложенные структуры и коллекции для изоляции слоёв.
func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepoModel(part.Category),
		Dimensions:    DimensionToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          copyStringSlice(part.Tags),
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// ┌───────────────────────────────────────────────────────┐
// │            Конвертеры для вложенных типов             │
// └───────────────────────────────────────────────────────┘

// MetadataFromRepoModel конвертирует метаданные из репозиторной модели в доменную.
func MetadataFromRepoModel(metadata map[string]repoModel.Value) map[string]model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]model.Value, len(metadata))
	for k, v := range metadata {
		result[k] = ValueFromRepoModel(v)
	}
	return result
}

// MetadataToRepoModel конвертирует метаданные из доменной модели в репозиторную.
func MetadataToRepoModel(metadata map[string]model.Value) map[string]repoModel.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]repoModel.Value, len(metadata))
	for k, v := range metadata {
		result[k] = ValueToRepoModel(v)
	}
	return result
}

// ValueFromRepoModel конвертирует одно значение метаданных (репозиторий → домен).
func ValueFromRepoModel(v repoModel.Value) model.Value {
	return model.Value{
		Kind:         model.ValueKind(v.Kind),
		StringValue:  v.StringValue,
		IntegerValue: v.IntegerValue,
		DoubleValue:  v.DoubleValue,
		BooleanValue: v.BooleanValue,
	}
}

// ValueToRepoModel конвертирует одно значение метаданных (домен → репозиторий).
func ValueToRepoModel(v model.Value) repoModel.Value {
	return repoModel.Value{
		Kind:         repoModel.ValueKind(v.Kind),
		StringValue:  v.StringValue,
		IntegerValue: v.IntegerValue,
		DoubleValue:  v.DoubleValue,
		BooleanValue: v.BooleanValue,
	}
}

// ManufacturerFromRepoModel конвертирует производителя (репозиторий → домен).
func ManufacturerFromRepoModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// ManufacturerToRepoModel конвертирует производителя (домен → репозиторий).
func ManufacturerToRepoModel(m model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// DimensionsFromRepoModel конвертирует размеры (репозиторий → домен).
// Примечание: если в твоих моделях используется имя "Dimension" (в единственном числе),
// убедись, что это согласовано в обеих структурах.
func DimensionsFromRepoModel(d repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// DimensionToRepoModel конвертирует размеры (домен → репозиторий).
func DimensionToRepoModel(d model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// CategoryFromRepoModel конвертирует категорию (репозиторий → домен).
func CategoryFromRepoModel(c repoModel.Category) model.Category {
	return model.Category(c)
}

// CategoryToRepoModel конвертирует категорию (домен → репозиторий).
func CategoryToRepoModel(c model.Category) repoModel.Category {
	return repoModel.Category(c)
}

// ┌───────────────────────────────────────────────────────┐
// │             Вспомогательные функции                   │
// └───────────────────────────────────────────────────────┘

// copyStringSlice создает глубокую копию среза строк.
// Сохраняет семантику nil: если входной срез nil — возвращает nil.
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
