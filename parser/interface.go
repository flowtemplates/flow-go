package parser

import (
	"github.com/flowtemplates/flow-go/lexer"
)

func AstFromBytes(input []byte) (Ast, error) {
	tokens := lexer.TokensFromBytes(input)
	ast, err := newParser(tokens).parse()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

// func ChanFromString(input string) <-chan Node {
// 	l := newLexer(input)
// 	return l.tokens
// }
