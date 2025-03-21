package parser

import (
	"github.com/flowtemplates/flow-go/lexer"
	"github.com/flowtemplates/flow-go/token"
)

type parser struct {
	tokens  []token.Token
	pos     int
	current token.Token
}

func newParser(tokens []token.Token) *parser {
	p := &parser{
		tokens: tokens,
		pos:    -1,
	}

	p.next()
	return p
}

func (p *parser) Parse() ([]Node, error) {
	var nodes []Node
	for p.pos < len(p.tokens) {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if node != nil {
			nodes = append(nodes, node)
		} else {
			p.next()
		}
	}

	return nodes, nil
}

func AstFromString(input string) ([]Node, error) {
	tokens := lexer.TokensFromString(input)
	ast, err := newParser(tokens).Parse()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

// func ChanFromString(input string) <-chan Node {
// 	l := newLexer(input)
// 	return l.tokens
// }
