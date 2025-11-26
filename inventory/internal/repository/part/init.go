package part

import (
	repoModel "github.com/clava1096/rocket-service/inventory/internal/repository/model"
)

func (r *repository) initSampleData() {
	engine := repoModel.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174000",
		Name:        "Main Engine",
		Description: "High-thrust rocket engine",
		Price:       150000.0,
		Category:    repoModel.CategoryEngine,
		Dimensions: repoModel.Dimensions{
			Length: 200.0,
			Width:  100.0,
			Height: 150.0,
			Weight: 500.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "CosmoEngines Inc.",
			Country: "USA",
			Website: "https://cosmoengines.example.com",
		},
		Tags: []string{"engine", "thrust", "rocket"},
	}
	r.inventory[engine.Uuid] = engine

	wing := repoModel.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Aerodynamic Wing",
		Description: "Lightweight composite wing for orbital stability",
		Price:       42000.50,
		Category:    repoModel.CategoryWing,
		Dimensions: repoModel.Dimensions{
			Length: 300.0,
			Width:  80.0,
			Height: 20.0,
			Weight: 120.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "AeroStructures Ltd.",
			Country: "Germany",
			Website: "https://aerostructures.example.com",
		},
		Tags: []string{"wing", "aero", "composite"},
	}
	r.inventory[wing.Uuid] = wing

	porthole := repoModel.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174002",
		Name:        "Reinforced Porthole",
		Description: "Triple-glazed viewport for deep-space observation",
		Price:       8500.75,
		Category:    repoModel.CategoryPortholes,
		Dimensions: repoModel.Dimensions{
			Length: 50.0,
			Width:  50.0,
			Height: 10.0,
			Weight: 25.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "ViewSpace Optics",
			Country: "France",
			Website: "https://viewspace.example.com",
		},
		Tags: []string{"window", "observation", "safe"},
	}
	r.inventory[porthole.Uuid] = porthole
}
