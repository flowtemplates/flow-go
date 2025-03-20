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

func render(ast []parser.Node, indent string, context Context) (string, error) {
	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case parser.Text:
			for _, s := range n.Val {
				result.WriteString(strings.TrimPrefix(s, indent))
			}
		case parser.ExprBlock:
			s, err := exprToValue(n.Body, context)
			if err != nil {
				return "", err
			}

			result.WriteString(s.String())
		case parser.IfStmt:
			conditionValue, err := exprToValue(n.BegTag.Body, context)
			if err != nil {
				return "", err
			}

			indent += n.BegTag.PreWs
			if conditionValue.Boolean() {
				bodyContent, err := render(n.Body, indent, context)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else if n.Else != nil {
				elseContent, err := render(n.Else, indent, context)
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
	case parser.Ident:
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
	case *parser.UnaryExpr:
		v, err := exprToValue(n.X, context)
		if err != nil {
			return nil, err
		}

		switch n.Op {
		case token.EXCL, token.NOT:
			return value.BooleanValue(!v.Boolean()), nil
		default:
			return nil, errors.New("unknown operator in unary expression")
		}
	case parser.Lit:
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

		switch n.Op {
		// case token.ADD:
		// 	return x.Add(y), nil
		case token.NEQL, token.ISNOT:
			return value.BooleanValue(x.String() != y.String()), nil
		case token.EQL, token.IS:
			return value.BooleanValue(x.String() == y.String()), nil
		case token.LAND, token.AND:
			if !x.Boolean() {
				return x, nil
			}
			return y, nil
		case token.LOR, token.OR:
			if x.Boolean() {
				return x, nil
			}
			return y, nil
		case token.GTR:
			return value.BooleanValue(x.Number() > y.Number()), nil
		case token.LESS:
			return value.BooleanValue(x.Number() < y.Number()), nil
		case token.LEQ:
			return value.BooleanValue(x.Number() <= y.Number()), nil
		case token.GEQ:
			return value.BooleanValue(x.Number() >= y.Number()), nil
		default:
			return nil, errors.New("unknown operator in binary expression")
		}
	default:
		return nil, fmt.Errorf("unsupported condition type: %T", expr)
	}
}
