package analyzer

import (
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/types"
)

func Typecheck(scope map[string]string, tm TypeMap) []TypeError {
	errs := []TypeError{}
	for name, typ := range tm {
		if typ == types.Any {
			continue
		}

		value, ok := scope[name]
		if !ok {
			scope[name] = typ.GetDefaultValue()
		} else if !typ.IsValid(value) {
			errs = append(errs, TypeError{
				ExpectedType: typ,
				Name:         name,
				Val:          value,
			})
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func GetTypeMapFromAst(ast []parser.Node, tm TypeMap) []error {
	errs := []error{}
	for _, node := range ast {
		switch n := node.(type) {
		case *parser.ExprNode:
			parseExpressionTypes(n.Body, tm, &errs)
		case *parser.IfNode:
			switch e := n.IfTag.Expr.(type) {
			case *parser.Ident:
				if err := addToTypeMap(Symbol{
					Name: e.Name,
					Typ:  types.Boolean,
				}, tm); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func GetTypeMapFromBytes(input []byte, tm TypeMap) error {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return err
	}

	if errs := GetTypeMapFromAst(ast, tm); len(errs) != 0 {
		return errs[0]
	}

	return nil
}
