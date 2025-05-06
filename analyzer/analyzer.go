package analyzer

import (
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

func (a analyzer) equalTypes(typ1, typ2 types.Type) types.Type {
	if xVarType, ok1 := typ1.(types.VarType); ok1 {
		t, err := a.Tm.addToTypeMap(string(xVarType), typ2)
		if err != nil {
			a.Errs.Add(err)
		}

		return t
	}

	return nil
}

func (a analyzer) parseNodes(ast parser.AST) {
	for _, node := range ast {
		switch n := node.(type) {
		case *parser.Print:
			a.parseExpressionTypes(n.Expr, types.String)

		case *parser.If:
			for _, condition := range n.Conditions {
				a.parseExpressionTypes(condition.Condition, types.Boolean)
				a.parseNodes(condition.Body)
			}

			a.parseNodes(n.ElseBody)

		case *parser.Switch:
			switchType := a.parseExpressionTypes(n.Expr, types.Any)

			for _, c := range n.Cases {
				caseTyp := a.parseExpressionTypes(c.Match, switchType)
				a.equalTypes(switchType, caseTyp)
				a.parseNodes(c.Body)
			}

			if n.Default != nil {
				a.parseNodes(n.Default)
			}
		}
	}
}

func (a analyzer) parseExpressionTypes(expr parser.Expr, typ types.Type) types.Type {
	switch e := expr.(type) {
	case *parser.Ident:
		t, err := a.Tm.addToTypeMap(e.Name, typ)
		if err != nil {
			a.Errs.Add(err)
		}

		return t

	// case *parser.BasicLit:
	// return e.Value.Type()

	case *parser.PipeExpr:
		a.parseExpressionTypes(e.X, types.String)

	case *parser.TernaryExpr:
		a.parseExpressionTypes(e.Condition, types.Boolean)
		a.parseExpressionTypes(e.ThenExpr, typ)
		a.parseExpressionTypes(e.ElseExpr, typ)

	case *parser.ParenExpr:
		a.parseExpressionTypes(e.X, typ)

	case *parser.BinaryExpr:
		if e.Op.IsLogicalOp() {
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
