package parser

import "github.com/flowtemplates/flow-go/lexer"

func AstFromString(input string) ([]Node, error) {
	tokens := lexer.TokensFromString(input)
	ast, errs := New(tokens).Parse()
	if len(errs) != 0 {
		return nil, errs[0]
	}

	return ast, nil
}
