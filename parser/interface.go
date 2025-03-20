package parser

import "github.com/flowtemplates/flow-go/lexer"

func AstFromString(input string) ([]Node, error) {
	tokens := lexer.TokensFromString(input)
	ast, err := New(tokens).Parse()
	if err != nil {
		return nil, err
	}

	return ast, nil
}
