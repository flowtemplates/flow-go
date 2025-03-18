package renderer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
)

type Scope map[string]any

type Context map[string]Valueable

func scopeToContext(scope Scope) Context {
	context := make(Context)
	for name, value := range scope {
		context[name] = ValueFromAny(value)
	}

	context["true"] = BooleanValue(true)
	context["false"] = BooleanValue(false)

	return context
}

func exprToValue(expr parser.Expr, context Context) (Valueable, error) {
	switch n := expr.(type) {
	case parser.Ident:
		value, exists := context[n.Name]
		if !exists {
			return nil, fmt.Errorf("%s not declared", n.Name)
		}

		return value, nil
	case parser.TernaryExpr:
		conditionValue, err := exprToValue(n.Condition, context)
		if err != nil {
			return nil, err
		}

		var exp parser.Expr
		if conditionValue.Boolean() {
			exp = n.TrueExpr
		} else {
			exp = n.FalseExpr
		}

		value, err := exprToValue(exp, context)
		if err != nil {
			return nil, err
		}

		return value, nil

	case parser.Lit:
		switch n.Typ {
		case token.STR:
			return StringValue(n.Val), nil
		case token.INT, token.FLOAT:
			// TODO:
			num, err := strconv.ParseFloat(n.Val, 64)
			if err != nil {
				return nil, err
			}

			return NumberValue(num), nil
		default:
			return nil, nil
		}

	case parser.BinaryExpr:
		x, err := exprToValue(n.X, context)
		if err != nil {
			return nil, err
		}

		y, err := exprToValue(n.Y, context)
		if err != nil {
			return nil, err
		}

		switch n.Op {
		case token.ADD:
			return x.Add(y), nil
		case token.EQL:
			return BooleanValue(x.String() == y.String()), nil
		default:
			return nil, errors.New("unknown operator in binary expression")
		}

	default:
		return nil, fmt.Errorf("unsupported condition type: %T", expr)
	}
}
