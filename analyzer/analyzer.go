package analyzer

import (
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/types"
)

type TypeMap map[string]types.Type

func (tm TypeMap) addToTypeMap(name string, typ types.Type) (types.Type, *TypeError) {
	if name == "true" || name == "false" {
		return types.Boolean, nil
	}

	if t, exists := tm[name]; !exists || t == types.Any {
		tm[name] = typ
	} else if typ != types.Any && t != typ {
		return types.Any, &TypeError{
			ExpectedType: typ,
			Name:         name,
			// Val:          string(v.Typ),
		}
	}

	return typ, nil
}

func parseNodes(ast []parser.Node, tm TypeMap, context renderer.Context, errs *TypeErrors) {
	for _, node := range ast {
		switch n := node.(type) {
		case *parser.ExprNode:
			parseExpressionTypes(n.Body, tm, types.Any, errs)
		case *parser.IfNode:
			parseExpressionTypes(n.IfTag.Expr, tm, types.Boolean, errs)
			parseNodes(n.Main, tm, context, errs)

			for _, elseIf := range n.ElseIfs {
				parseExpressionTypes(elseIf.Tag.Expr, tm, types.Boolean, errs)
				parseNodes(elseIf.Body, tm, context, errs)
			}

			parseExpressionTypes(n.IfTag.Expr, tm, types.Boolean, errs)
			parseNodes(n.Main, tm, context, errs)
		}
	}
}

func parseExpressionTypes(expr parser.Expr, tm TypeMap, typ types.Type, errs *TypeErrors) types.Type {
	switch e := expr.(type) {
	case *parser.Ident:
		t, err := tm.addToTypeMap(e.Name, typ)
		if err != nil {
			errs.Add(err)
		}

		return t
	case *parser.StringLit:
		return e.Value.Type()
	case *parser.NumberLit:
		return e.Value.Type()
	case *parser.FilterExpr:
		parseExpressionTypes(e.Expr, tm, types.String, errs)
	case *parser.TernaryExpr:
		parseExpressionTypes(e.Condition, tm, types.Boolean, errs)
		parseExpressionTypes(e.TrueExpr, tm, types.Any, errs)
		parseExpressionTypes(e.FalseExpr, tm, types.Any, errs)
	case *parser.BinaryExpr:
		switch e.Op.Kind {
		case token.GRTR, token.LESS, token.LEQ, token.GEQ:
			parseExpressionTypes(e.X, tm, types.Number, errs)
			parseExpressionTypes(e.Y, tm, types.Number, errs)
		default:
			parseExpressionTypes(e.X, tm, types.Any, errs)
			parseExpressionTypes(e.Y, tm, types.Any, errs)
		}

		return types.Any
	}

	return types.Any
}
