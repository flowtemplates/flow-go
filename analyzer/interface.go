package analyzer

import (
	"github.com/flowtemplates/flow-go/lexer"
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
		case parser.ExprBlock:
			parseExpressionTypes(n.Body, tm, &errs)
		case parser.IfStmt:
			switch e := n.BegTag.Body.(type) {
			case parser.Ident:
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

func GetTypeMapFromString(input string, tm TypeMap) error {
	tokens := lexer.TokensFromString(input)
	ast, errs := parser.New(tokens).Parse()
	if len(errs) != 0 {
		return errs[0]
	}

	if errs := GetTypeMapFromAst(ast, tm); len(errs) != 0 {
		return errs[0]
	}

	return nil
}
