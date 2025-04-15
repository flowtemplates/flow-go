package renderer

import (
	"bytes"
	"errors"
	"fmt"

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

	// TODO: check overwrite
	context["true"] = value.BooleanValue(true)
	context["false"] = value.BooleanValue(false)

	return context
}

func render(ast []parser.Node, context Context) ([]byte, error) {
	var buf bytes.Buffer

	for _, node := range ast {
		switch n := node.(type) {
		case *parser.TextNode:
			for _, s := range n.Val {
				buf.WriteString(s)
			}

		case *parser.ExprNode:
			s, err := exprToValue(n.Body, context)
			if err != nil {
				return nil, err
			}

			buf.WriteString(s.AsString())

		case *parser.IfNode:
			conditionValue, err := exprToValue(n.IfTag.Expr, context)
			if err != nil {
				return nil, err
			}

			if conditionValue.AsBoolean() {
				bodyContent, err := render(n.Main, context)
				if err != nil {
					return nil, err
				}

				buf.Write(bodyContent)

				return buf.Bytes(), nil
			}

			for _, elseIf := range n.ElseIfs {
				elifCondition, err := exprToValue(elseIf.Tag.Expr, context)
				if err != nil {
					return nil, err
				}

				if elifCondition.AsBoolean() {
					elifContent, err := render(elseIf.Body, context)
					if err != nil {
						return nil, err
					}

					buf.Write(elifContent)

					return buf.Bytes(), nil
				}
			}

			elseContent, err := render(n.Else.Body, context)
			if err != nil {
				return nil, err
			}

			buf.Write(elseContent)

		case *parser.SwitchNode:
			switchValue, err := exprToValue(n.SwitchTag.Expr, context)
			if err != nil {
				return nil, err
			}

			caseMatched := false

			for _, c := range n.Cases {
				val, err := exprToValue(c.Tag.Expr, context)
				if err != nil {
					return nil, err
				}

				if eql(switchValue, val) {
					body, err := render(c.Body, context)
					if err != nil {
						return nil, err
					}

					buf.Write(body)

					caseMatched = true

					break
				}
			}

			if !caseMatched && n.DefaultCase != nil {
				body, err := render(n.DefaultCase.Body, context)
				if err != nil {
					return nil, err
				}

				buf.Write(body)
			}

		default:
			return nil, fmt.Errorf("unexpected node type in ast: %T", n)
		}
	}

	return buf.Bytes(), nil
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

	case *parser.FilterExpr:
		expr, err := exprToValue(n.Expr, context)
		if err != nil {
			return nil, err
		}

		return callFilter(n.Filter.Name, expr)

	case *parser.ParenExpr:
		return exprToValue(n.Expr, context)

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
