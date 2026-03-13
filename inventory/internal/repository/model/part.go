package model

import (
	"time"
)

type Part struct {
	Uuid          string           `bson:"_id,omitempty"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float64          `bson:"price"`
	StockQuantity float64          `bson:"stock_quantity"`
	Category      Category         `bson:"category"`
	Dimensions    Dimensions       `bson:"dimensions"`
	Manufacturer  Manufacturer     `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     *time.Time       `bson:"updated_at,omitempty"`
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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type PartsFilter struct {
	Uuids                 []string   `bson:"uuids"`
	Names                 []string   `bson:"names"`
	Categories            []Category `bson:"categories"`
	ManufacturerCountries []string   `bson:"manufacturer_countries"`
	Tags                  []string   `bson:"tags"`
}

type Value struct {
	Kind         ValueKind `bson:"kind"`
	StringValue  string    `bson:"string_value"`
	IntegerValue int64     `bson:"integer_value"`
	DoubleValue  float64   `bson:"double_value"`
	BooleanValue bool      `bson:"boolean_value"`
}

type ValueKind int

const (
	ValueKindUnspecified ValueKind = iota
	ValueKindString
	ValueKindInt64
	ValueKindFloat64
	ValueKindBool
)
