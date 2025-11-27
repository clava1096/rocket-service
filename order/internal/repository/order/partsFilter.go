package model

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Category string

const (
	CategoryUnspecified Category = ""
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryWing        Category = "WING"
	CategoryPortholes   Category = "PORTHOLE"
)
