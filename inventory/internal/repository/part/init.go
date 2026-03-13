package part

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/clava1096/rocket-service/inventory/internal/model"
)

func (r *repository) initSampleData() {
	ctx := context.Background()
	now := time.Now()

	parts := []model.Part{
		{
			Uuid:          uuid.New().String(),
			Name:          "Main Engine",
			Description:   "High-thrust rocket engine",
			Price:         150000.0,
			StockQuantity: 10,
			Category:      model.CategoryEngine,
			Dimensions: model.Dimensions{
				Length: 200.0,
				Width:  100.0,
				Height: 150.0,
				Weight: 500.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "CosmoEngines Inc.",
				Country: "USA",
				Website: "https://cosmoengines.example.com",
			},
			Tags:      []string{"engine", "thrust", "rocket"},
			CreatedAt: now,
		},
		{
			Uuid:          uuid.New().String(),
			Name:          "Aerodynamic Wing",
			Description:   "Lightweight composite wing for orbital stability",
			Price:         42000.50,
			StockQuantity: 15,
			Category:      model.CategoryWing,
			Dimensions: model.Dimensions{
				Length: 300.0,
				Width:  80.0,
				Height: 20.0,
				Weight: 120.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "AeroStructures Ltd.",
				Country: "Germany",
				Website: "https://aerostructures.example.com",
			},
			Tags:      []string{"wing", "aero", "composite"},
			CreatedAt: now,
		},
		{
			Uuid:          uuid.New().String(),
			Name:          "Reinforced Porthole",
			Description:   "Triple-glazed viewport for deep-space observation",
			Price:         8500.75,
			StockQuantity: 30,
			Category:      model.CategoryPortholes,
			Dimensions: model.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "ViewSpace Optics",
				Country: "France",
				Website: "https://viewspace.example.com",
			},
			Tags:      []string{"window", "observation", "safe"},
			CreatedAt: now,
		},
	}

	for _, part := range parts {
		_, err := r.Create(ctx, part)
		if err != nil {
			log.Fatalf("Failed to create part '%s': %v", part.Name, err)
		}
	}
}
