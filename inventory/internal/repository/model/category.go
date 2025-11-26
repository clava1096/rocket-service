package model

type Category string

const (
	CategoryUnspecified Category = ""
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryWing        Category = "WING"
	CategoryPortholes   Category = "PORTHOLE"
)
