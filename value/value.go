package value

import (
	"fmt"
	"math"
	"strconv"

	"github.com/flowtemplates/flow-go/types"
)

type Valueable interface {
	String() string
	Boolean() bool
	Number() float64
	Add(Valueable) Valueable
	Type() types.Type
}

func FromAny(value any) Valueable {
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

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Boolean() bool {
	return v != ""
}

func (v StringValue) Number() float64 {
	if num, err := strconv.ParseFloat(string(v), 64); err == nil {
		return num
	}

	var sum float64
	for i, char := range v {
		sum += float64(char) * float64(i+1) // Weighted position for uniqueness
	}

	return sum
}

func (v StringValue) Add(b Valueable) Valueable {
	return StringValue(string(v) + b.String())
}

func (v StringValue) Type() types.Type {
	return types.String
}

type BooleanValue bool

func (v BooleanValue) String() string {
	return ""
}

func (v BooleanValue) Boolean() bool {
	return bool(v)
}

func (v BooleanValue) Number() float64 {
	if v {
		return 1
	}

	return 0
}

func (v BooleanValue) Add(b Valueable) Valueable {
	switch b.(type) {
	case BooleanValue, NumberValue:
		return NumberValue(v.Number() + b.Number())
	default:
		return StringValue(v.String() + b.String())
	}
}

func (v BooleanValue) Type() types.Type {
	return types.Boolean
}

type NumberValue float64

func (v NumberValue) String() string {
	floatValue := float64(v)
	if math.Trunc(floatValue) == floatValue {
		return fmt.Sprintf("%.0f", v)
	}

	return fmt.Sprintf("%g", v)
}

func (v NumberValue) Boolean() bool {
	return v != 0
}

func (v NumberValue) Number() float64 {
	return float64(v)
}

func (v NumberValue) Add(b Valueable) Valueable {
	switch b.(type) {
	case BooleanValue, NumberValue:
		return NumberValue(v.Number() + b.Number())
	default:
		return StringValue(v.String() + b.String())
	}
}

func (v NumberValue) Type() types.Type {
	return types.Number
}
