package renderer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

type Input map[string]any

// Map with variables and their values
type Context map[string]value.Valuable

func InputToContext(scope Input) Context {
	context := make(Context)
	for name, val := range scope {
		context[name] = value.FromAny(val)
	}

	return context
}

func render(ast parser.AST, context Context) (string, error) {
	var buf strings.Builder

	for _, node := range ast {
		switch n := node.(type) {
		case *parser.Text:
			buf.WriteString(n.Value)

		case *parser.Print:
			s, err := exprToValue(n.Expr, context)
			if err != nil {
				return "", err
			}

			buf.WriteString(s.AsString())

		case *parser.If:
			for _, elseIf := range n.Conditions {
				elifCondition, err := exprToValue(elseIf.Condition, context)
				if err != nil {
					return "", err
				}

				if elifCondition.AsBoolean() {
					elifContent, err := render(elseIf.Body, context)
					if err != nil {
						return "", err
					}

					buf.WriteString(elifContent)
				}
			}

			elseContent, err := render(n.ElseBody, context)
			if err != nil {
				return "", err
			}

			buf.WriteString(elseContent)

		case *parser.Switch:
			switchValue, err := exprToValue(n.Expr, context)
			if err != nil {
				return "", err
			}

			caseMatched := false

			for _, c := range n.Cases {
				val, err := exprToValue(c.Match, context)
				if err != nil {
					return "", err
				}

				if eql(switchValue, val) {
					body, err := render(c.Body, context)
					if err != nil {
						return "", err
					}

					buf.WriteString(body)

					caseMatched = true

					break
				}
			}

			if !caseMatched && n.Default != nil {
				body, err := render(n.Default, context)
				if err != nil {
					return "", err
				}

				buf.WriteString(body)
			}

		default:
			return "", fmt.Errorf("unexpected node type in ast: %T", n)
		}
	}

	return buf.String(), nil
}

func eql(x, y value.Valuable) bool {
	return x.AsString() == y.AsString()
}

func exprToValue(expr parser.Expr, context Context) (value.Valuable, error) {
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
			exp = n.ThenExpr
		} else {
			exp = n.ElseExpr
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

		switch n.Op {
		case token.EXCL, token.NOT:
			return value.BooleanValue(!v.AsBoolean()), nil

		default:
			return nil, errors.New("unknown operator in unary expression")
		}

	case *parser.BasicLit:
		return n.Value, nil

	case *parser.PipeExpr:
		expr, err := exprToValue(n.X, context)
		if err != nil {
			return nil, err
		}

		var val value.Valuable
		for _, f := range n.Filters {
			val, err = callFilter(f.Name, expr)
			if err != nil {
				return nil, err
			}
		}

		return val, nil

	case *parser.ParenExpr:
		return exprToValue(n.X, context)

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
		case token.ADD:
			return x.Add(y), nil

		case token.NEQL, token.ISNOT:
			return value.BooleanValue(x.AsString() != y.AsString()), nil

		case token.EQL, token.IS:
			return value.BooleanValue(eql(x, y)), nil

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
