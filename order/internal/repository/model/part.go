package model

import "time"

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity float64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnspecified Category = ""
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryWing        Category = "WING"
	CategoryPortholes   Category = "PORTHOLE"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Value struct {
	Kind         ValueKind
	StringValue  string
	IntegerValue int64
	DoubleValue  float64
	BooleanValue bool
}

type ValueKind int

const (
	ValueKindUnspecified ValueKind = iota
	ValueKindString
	ValueKindInt64
	ValueKindFloat64
	ValueKindBool
)
