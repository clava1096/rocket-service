// internal/service/part/list_test.go
package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	var parts []model.Part
	for i := 0; i < 3; i++ {
		part := model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.Phrase(),
			Price:         gofakeit.Float64Range(10, 10000),
			StockQuantity: gofakeit.Float64(),
			Category:      s.randomCategory(),
			Dimensions: model.Dimensions{
				Length: gofakeit.Float64Range(10, 200),
				Width:  gofakeit.Float64Range(5, 100),
				Height: gofakeit.Float64Range(5, 100),
				Weight: gofakeit.Float64Range(1, 500),
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.Hobby(), gofakeit.Adjective()},
			Metadata: map[string]model.Value{
				"color": {
					Kind:        model.ValueKindString,
					StringValue: gofakeit.Color(),
				},
			},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
		}
		parts = append(parts, part)
	}

	filter := model.PartsFilter{
		Uuids: []string{parts[0].Uuid}, // опционально
	}

	s.partRepository.On("List", s.ctx, filter).Return(parts, nil)

	actualParts, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.NotNil(actualParts)
	s.Len(actualParts, 3)

	s.Equal(parts[0].Uuid, actualParts[0].Uuid)
	s.Equal(parts[0].Name, actualParts[0].Name)
	s.Equal(parts[0].Price, actualParts[0].Price)
	s.Equal(parts[0].Category, actualParts[0].Category)
	s.Equal(parts[0].Dimensions, actualParts[0].Dimensions)
	s.Equal(parts[0].Manufacturer, actualParts[0].Manufacturer)
	s.Equal(parts[0].Tags, actualParts[0].Tags)
	s.Equal(parts[0].Metadata, actualParts[0].Metadata)
}

func (s *ServiceSuite) TestListRepoError() {
	repoError := gofakeit.Error()

	filter := model.PartsFilter{}

	s.partRepository.On("List", s.ctx, filter).Return([]model.Part{}, repoError)

	actualParts, err := s.service.List(s.ctx, filter)

	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Empty(actualParts)
}

func (s *ServiceSuite) randomCategory() model.Category {
	categories := []model.Category{
		model.CategoryEngine,
		model.CategoryFuel,
		model.CategoryPortholes,
		model.CategoryWing,
	}
	return categories[gofakeit.Number(0, len(categories)-1)]
}

func (s *ServiceSuite) TestGetNotFound() {
	uuid := gofakeit.UUID()
	s.partRepository.On("Get", s.ctx, uuid).Return(model.Part{}, model.ErrNotFound)

	part, err := s.service.Get(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
	s.Empty(part)
}
