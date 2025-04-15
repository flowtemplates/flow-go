package types

import "fmt"

type Type interface {
	t()
	// String() string
	// For unification
	// apply(subst map[string]Type) Type
	// ftv() map[string]bool
}

type PrimitiveType string

const (
	Number  PrimitiveType = "number"
	String  PrimitiveType = "string"
	Boolean PrimitiveType = "boolean"
	// TODO: move to separate struct or remove completely
	Any PrimitiveType = "any"
	// TypeArray   Type = "array"
	// TypeObject  Type = "object"
)

func (t PrimitiveType) t() {}

func (t PrimitiveType) IsValid(val any) bool {
	fmt.Printf("%T, %s\n", val, val)
	switch t {
	case Number:
		switch val.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			return true

		default:
			return false
		}

	case Boolean:
		_, ok := val.(bool)

		return ok

	case String:
		_, ok := val.(string)

		return ok

	case Any:
		return true

	default:
		return false
	}
}

func (t PrimitiveType) GetDefaultValue() string {
	switch t {
	case Number:
		return "0"

	case Boolean:
		return "false"

	default:
		return ""
	}
}

type VarType string

func (t VarType) t() {}
