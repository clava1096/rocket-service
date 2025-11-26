package model

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
