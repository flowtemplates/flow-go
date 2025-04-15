package analyzer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/types"
)

type TypeMap map[string]types.Type

func (tm TypeMap) getPrimitive(typ types.Type) *types.PrimitiveType {
	if primmitiveType, ok := typ.(types.PrimitiveType); ok {
		return &primmitiveType
	}

	if varType, ok := typ.(types.VarType); ok {
		t, ok := tm[string(varType)]
		if !ok {
			return nil
		}

		if primmitiveType, ok := t.(types.PrimitiveType); ok {
			return &primmitiveType
		}

		fmt.Println("recursive get primitive")

		return tm.getPrimitive(t)
	}

	return nil
}

func (tm TypeMap) addToTypeMap(name string, typ types.Type) (types.Type, *TypeError) {
	// FIXME
	if name == "true" || name == "false" {
		return types.Boolean, nil
	}

	current, exists := tm[name]

	if varType, ok := current.(types.VarType); ok {
		return tm.addToTypeMap(string(varType), typ)
	}

	switch {
	case !exists || current == types.Any:
		tm[name] = typ

	case current == types.Boolean:
		if typ != types.Any && typ != types.Boolean {
			tm[name] = typ
		}

	case typ != types.Any && current != typ:
		return nil, &TypeError{
			Name:         name,
			ExpectedType: *tm.getPrimitive(typ),
		}
	}

	return types.VarType(name), nil
}

func (a *Analyzer) equalTypes(typ1, typ2 types.Type) types.Type {
	if xVarType, ok1 := typ1.(types.VarType); ok1 {
		t, err := a.Tm.addToTypeMap(string(xVarType), typ2)
		if err != nil {
			a.Errs.Add(err)
		}

		return t
	}

	return nil
}

func (a *Analyzer) parseNodes(ast []parser.Node) {
	for _, node := range ast {
		switch n := node.(type) {
		case *parser.ExprNode:
			a.parseExpressionTypes(n.Body, types.String)

		case *parser.IfNode:
			a.parseExpressionTypes(n.IfTag.Expr, types.Boolean)
			a.parseNodes(n.Main)

			for _, elseIf := range n.ElseIfs {
				a.parseExpressionTypes(elseIf.Tag.Expr, types.Boolean)
				a.parseNodes(elseIf.Body)
			}

			a.parseExpressionTypes(n.IfTag.Expr, types.Boolean)
			a.parseNodes(n.Main)

		case *parser.SwitchNode:
			switchType := a.parseExpressionTypes(n.SwitchTag.Expr, types.Any)

			for _, c := range n.Cases {
				caseTyp := a.parseExpressionTypes(c.Tag.Expr, switchType)
				a.equalTypes(switchType, caseTyp)
				a.parseNodes(c.Body)
			}

			if n.DefaultCase != nil {
				a.parseNodes(n.DefaultCase.Body)
			}
		}
	}
}

func (a *Analyzer) parseExpressionTypes(expr parser.Expr, typ types.Type) types.Type {
	switch e := expr.(type) {
	case *parser.Ident:
		t, err := a.Tm.addToTypeMap(e.Name, typ)
		if err != nil {
			a.Errs.Add(err)
		}

		return t

	case *parser.StringLit:
		return e.Value.Type()

	case *parser.NumberLit:
		return e.Value.Type()

	case *parser.FilterExpr:
		a.parseExpressionTypes(e.Expr, types.String)

	case *parser.TernaryExpr:
		a.parseExpressionTypes(e.Condition, types.Boolean)
		a.parseExpressionTypes(e.TrueExpr, typ)
		a.parseExpressionTypes(e.FalseExpr, typ)

	case *parser.ParenExpr:
		a.parseExpressionTypes(e.Expr, typ)

	case *parser.BinaryExpr:
		if e.Op.Kind.IsLogicalOp() {
			a.parseExpressionTypes(e.X, types.Boolean)
			a.parseExpressionTypes(e.Y, types.Boolean)
		} else {
			x := a.parseExpressionTypes(e.X, types.Any)
			y := a.parseExpressionTypes(e.Y, types.Any)

			return a.equalTypes(x, y)
		}
	}

	return types.Any
}
