package analyzer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/types"
)

type Symbol struct {
	Name string
	Typ  types.Type
}

type TypeMap map[string]types.Type

func addToTypeMap(v Symbol, tm TypeMap) error {
	if typ, exists := tm[v.Name]; !exists || typ == types.Any {
		tm[v.Name] = v.Typ
	} else if v.Typ != types.Any && v.Typ != typ {
		return fmt.Errorf("unmatched type of %q, got %s, expected %s", v.Name, v.Typ, typ)
	}

	return nil
}

func parseExpressionTypes(expr parser.Expr, tm TypeMap, errs *[]error) types.Type {
	switch e := expr.(type) {
	case parser.Ident:
		if e.Name == "true" || e.Name == "false" {
			return types.Boolean
		}

		if err := addToTypeMap(Symbol{
			Name: e.Name,
			Typ:  types.Any,
		}, tm); err != nil {
			*errs = append(*errs, err)
		}

		return types.Any
	case parser.Lit:
		return e.Value.Type()
	case parser.TernaryExpr:
		parseExpressionTypes(e.Condition, tm, errs)
		parseExpressionTypes(e.TrueExpr, tm, errs)
		parseExpressionTypes(e.FalseExpr, tm, errs)
	case parser.BinaryExpr:
		// t1 := parseExpressionTypes(e.X, tm, errs)
		// t2 := parseExpressionTypes(e.Y, tm, errs)

		// DO NOT DELETE JUST IN CASE
		// if e.Op == token.ADD {
		// 	if t1 == types.String || t2 == types.String {
		// 		// If one side is a string, enforce string type
		// 		if ident, ok := e.X.(parser.Ident); ok {
		// 			if err := addToTypeMap(Symbol{Name: ident.Name, Typ: types.String}, tm); err != nil {
		// 				*errs = append(*errs, err)
		// 			}
		// 		}
		// 		if ident, ok := e.Y.(parser.Ident); ok {
		// 			if err := addToTypeMap(Symbol{Name: ident.Name, Typ: types.String}, tm); err != nil {
		// 				*errs = append(*errs, err)
		// 			}
		// 		}
		// 		return types.String
		// 	} else if t1 == types.Number || t2 == types.Number {
		// 		// If one side is a number, enforce number type
		// 		if ident, ok := e.X.(parser.Ident); ok {
		// 			if err := addToTypeMap(Symbol{Name: ident.Name, Typ: types.Number}, tm); err != nil {
		// 				*errs = append(*errs, err)
		// 			}
		// 		}
		// 		if ident, ok := e.Y.(parser.Ident); ok {
		// 			if err := addToTypeMap(Symbol{Name: ident.Name, Typ: types.Number}, tm); err != nil {
		// 				*errs = append(*errs, err)
		// 			}
		// 		}

		// 		return types.Number
		// 	}
		// }

		return types.Any
	}

	return types.Any
}
