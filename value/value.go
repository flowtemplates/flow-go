package value

import (
	"fmt"
	"math"
	"strconv"

	"github.com/flowtemplates/flow-go/types"
)

type Valuable interface {
	AsString() string
	AsBoolean() bool
	AsNumber() float64
	Add(value Valuable) Valuable
	Type() types.PrimitiveType
}

func FromAny(value any) Valuable {
	switch v := value.(type) {
	case string:
		return StringValue(v)

	case *string:
		return StringValue(*v)

	case float64:
		return NumberValue(v)

	case *float64:
		return NumberValue(*v)

	case int:
		return NumberValue(v)

	case *int:
		return NumberValue(*v)

	case bool:
		return BooleanValue(v)

	case *bool:
		return BooleanValue(*v)

	default:
		panic(fmt.Sprintf("cannot convert any to Valuable: unsupported type: %T", value))
	}
}

type StringValue string

func (v StringValue) AsString() string {
	return string(v)
}

func (v StringValue) AsBoolean() bool {
	return v != ""
}

func (v StringValue) AsNumber() float64 {
	if num, err := strconv.ParseFloat(string(v), 64); err == nil {
		return num
	}

	var sum float64
	for i, char := range v {
		sum += float64(char) * float64(i+1) // Weighted position for uniqueness
	}

	return sum
}

func (v StringValue) Add(b Valuable) Valuable {
	return StringValue(string(v) + b.AsString())
}

func (v StringValue) Type() types.PrimitiveType {
	return types.String
}

type BooleanValue bool

func (v BooleanValue) AsString() string {
	return ""
}

func (v BooleanValue) AsBoolean() bool {
	return bool(v)
}

func (v BooleanValue) AsNumber() float64 {
	if v {
		return 1
	}

	return 0
}

func (v BooleanValue) Add(b Valuable) Valuable {
	switch b.(type) {
	case BooleanValue, NumberValue:
		return NumberValue(v.AsNumber() + b.AsNumber())

	default:
		return StringValue(v.AsString() + b.AsString())
	}
}

func (v BooleanValue) Type() types.PrimitiveType {
	return types.Boolean
}

type NumberValue float64

func (v NumberValue) AsString() string {
	floatValue := float64(v)
	if math.Trunc(floatValue) == floatValue {
		return fmt.Sprintf("%.0f", v)
	}

	return fmt.Sprintf("%g", v)
}

func (v NumberValue) AsBoolean() bool {
	return v != 0
}

func (v NumberValue) AsNumber() float64 {
	return float64(v)
}

func (v NumberValue) Add(b Valuable) Valuable {
	switch b.(type) {
	case BooleanValue, NumberValue:
		return NumberValue(v.AsNumber() + b.AsNumber())

	default:
		return StringValue(v.AsString() + b.AsString())
	}
}

func (v NumberValue) Type() types.PrimitiveType {
	return types.Number
}
