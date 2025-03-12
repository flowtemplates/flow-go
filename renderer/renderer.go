package renderer

import (
	"fmt"
	"strconv"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/types"
)

type Scope map[string]string

type Context map[string]TypedValue

type TypedValue struct {
	Typ types.Type
	Val any
}

func (s TypedValue) String() string {
	switch s.Typ {
	case types.Boolean:
		return ""
	case types.Number:
		return fmt.Sprintf("%d", s.Val)
	default:
		return fmt.Sprintf("%s", s.Val)
	}
}

func (s TypedValue) Boolean() bool {
	switch s.Typ {
	case types.Boolean:
		return s.Val.(bool)
	case types.Number:
		return s.Val.(float64) != 0
	case types.String:
		return s.Val.(string) != ""
	default:
		return true
	}
}

func scopeToContext(scope Scope, tm analyzer.TypeMap) Context {
	context := make(Context)
	for name, value := range scope {
		typeInfo, exists := tm[name]
		if !exists {
			continue // Skip unknown types
		}

		var typedValue any
		switch typeInfo {
		case types.Number:
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				typedValue = v
			} else {
				typedValue = 0.0 // Default if conversion fails
			}
		case types.Boolean:
			typedValue = value == "true"
		default:
			typedValue = value // Treat as string or unknown type
		}

		context[name] = TypedValue{
			Typ: typeInfo,
			Val: typedValue,
		}
	}

	context["true"] = TypedValue{
		Typ: types.Boolean,
		Val: true,
	}
	context["false"] = TypedValue{
		Typ: types.Boolean,
		Val: false,
	}

	return context
}

func exprToValue(cond parser.Expr, context Context) (TypedValue, error) {
	switch n := cond.(type) {
	case parser.Ident:
		value, exists := context[n.Name]
		if !exists {
			return TypedValue{}, fmt.Errorf("%s not declared", n.Name)
		}

		return value, nil
	default:
		return TypedValue{}, fmt.Errorf("unsupported condition type: %T", cond)
	}
}
