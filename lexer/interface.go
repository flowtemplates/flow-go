package lexer

import (
	"github.com/flowtemplates/flow-go/token"
)

type lexer struct {
	source   []byte
	startPos token.Position
	pos      token.Position
	tokensCh chan token.Token
}

func newLexer(input []byte) *lexer {
	l := lexer{
		source: input,
		startPos: token.Position{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
		tokensCh: make(chan token.Token, 2),
	}
	l.pos = l.startPos

	go l.run()
	return &l
}

func TokensFromBytes(source []byte) []token.Token {
	l := newLexer(source)

	var tokens []token.Token
	for {
		tok := <-l.tokensCh
		if tok.Kind == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens
}

func ChanFromBytes(source []byte) <-chan token.Token {
	l := newLexer(source)
	return l.tokensCh
}
