package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		uuid          = gofakeit.UUID()
		name          = gofakeit.Name()
		desc          = gofakeit.Phrase()
		price         = gofakeit.Float64()
		stockQuantity = gofakeit.Float64()
		length        = gofakeit.Float64()
		width         = gofakeit.Float64()
		height        = gofakeit.Float64()
		weight        = gofakeit.Float64()
		tags          = []string{gofakeit.Word(), gofakeit.Word()}
		manufacturer  = model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
		metadata = map[string]model.Value{
			"color": {
				Kind:        model.ValueKindString,
				StringValue: "red",
			},
			"weight_kg": {
				Kind:        model.ValueKindFloat64,
				DoubleValue: weight,
			},
		}
	)

	now := time.Now().Truncate(time.Second)
	expectedPartValue := model.Part{
		Uuid:          uuid,
		Name:          name,
		Description:   desc,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      model.CategoryPortholes,
		Dimensions: model.Dimensions{
			Length: length,
			Width:  width,
			Height: height,
			Weight: weight,
		},
		Manufacturer: manufacturer,
		Tags:         tags,
		Metadata:     metadata,
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	s.partRepository.On("Get", s.ctx, uuid).Return(expectedPartValue, nil)

	actualPart, err := s.service.Get(s.ctx, uuid)

	s.NoError(err)
	s.NotNil(actualPart)
	s.Equal(expectedPartValue.Uuid, actualPart.Uuid)
	s.Equal(expectedPartValue.Name, actualPart.Name)
	s.Equal(expectedPartValue.Description, actualPart.Description)
	s.Equal(expectedPartValue.Price, actualPart.Price)
	s.Equal(expectedPartValue.StockQuantity, actualPart.StockQuantity)
	s.Equal(expectedPartValue.Category, actualPart.Category)
	s.Equal(expectedPartValue.Dimensions, actualPart.Dimensions)
	s.Equal(expectedPartValue.Manufacturer, actualPart.Manufacturer)
	s.Equal(expectedPartValue.Tags, actualPart.Tags)
	s.Equal(expectedPartValue.Metadata, actualPart.Metadata)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoError = gofakeit.Error()
		uuid      = gofakeit.UUID()
	)

	s.partRepository.On("Get", s.ctx, uuid).Return(model.Part{}, repoError)

	actualPart, err := s.service.Get(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoError)
	s.Empty(actualPart)
}
