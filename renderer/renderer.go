package renderer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

type Scope map[string]any

type Context map[string]value.Valueable

func scopeToContext(scope Scope) Context {
	context := make(Context)
	for name, val := range scope {
		context[name] = value.FromAny(val)
	}

	context["true"] = value.BooleanValue(true)
	context["false"] = value.BooleanValue(false)

	return context
}

func render(ast []parser.Node, context Context) (string, error) {
	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case *parser.TextNode:
			for _, s := range n.Val {
				result.WriteString(s)
			}
		case *parser.ExprNode:
			s, err := exprToValue(n.Body, context)
			if err != nil {
				return "", err
			}

			result.WriteString(s.AsString())
		case *parser.IfNode:
			conditionValue, err := exprToValue(n.IfTag.Expr, context)
			if err != nil {
				return "", err
			}

			if conditionValue.AsBoolean() {
				bodyContent, err := render(n.MainBody, context)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else {
				elseContent, err := render(n.ElseBody.Body, context)
				if err != nil {
					return "", err
				}

				result.WriteString(elseContent)
			}
		}
	}

	return result.String(), nil
}

func exprToValue(expr parser.Expr, context Context) (value.Valueable, error) {
	switch n := expr.(type) {
	case *parser.Ident:
		value, exists := context[n.Name]
		if !exists {
			return nil, fmt.Errorf("%s not declared", n.Name)
		}

		return value, nil
	case *parser.TernaryExpr:
		conditionValue, err := exprToValue(n.Condition, context)
		if err != nil {
			return nil, err
		}

		var exp parser.Expr
		if conditionValue.AsBoolean() {
			exp = n.TrueExpr
		} else {
			exp = n.FalseExpr
		}

		value, err := exprToValue(exp, context)
		if err != nil {
			return nil, err
		}

		return value, nil
	case *parser.UnaryExpr:
		v, err := exprToValue(n.Expr, context)
		if err != nil {
			return nil, err
		}

		switch n.Op.Kind {
		case token.EXCL, token.NOT:
			return value.BooleanValue(!v.AsBoolean()), nil
		default:
			return nil, errors.New("unknown operator in unary expression")
		}
	case *parser.NumberLit:
		return n.Value, nil
	case *parser.StringLit:
		return n.Value, nil
	case *parser.BinaryExpr:
		x, err := exprToValue(n.X, context)
		if err != nil {
			return nil, err
		}

		y, err := exprToValue(n.Y, context)
		if err != nil {
			return nil, err
		}

		switch n.Op.Kind {
		// case token.ADD:
		// 	return x.Add(y), nil
		case token.NEQL, token.ISNOT:
			return value.BooleanValue(x.AsString() != y.AsString()), nil
		case token.EQL, token.IS:
			return value.BooleanValue(x.AsString() == y.AsString()), nil
		case token.LAND, token.AND:
			if !x.AsBoolean() {
				return x, nil
			}
			return y, nil
		case token.LOR, token.OR:
			if x.AsBoolean() {
				return x, nil
			}
			return y, nil
		case token.GRTR:
			return value.BooleanValue(x.AsNumber() > y.AsNumber()), nil
		case token.LESS:
			return value.BooleanValue(x.AsNumber() < y.AsNumber()), nil
		case token.LEQ:
			return value.BooleanValue(x.AsNumber() <= y.AsNumber()), nil
		case token.GEQ:
			return value.BooleanValue(x.AsNumber() >= y.AsNumber()), nil
		default:
			return nil, errors.New("unknown operator in binary expression")
		}
	default:
		return nil, fmt.Errorf("unsupported condition type: %T", expr)
	}
}
