package lexer

import (
	"github.com/flowtemplates/flow-go/token"
)

type lexer struct {
	input  string
	start  token.Position
	pos    token.Position
	tokens chan token.Token
}

func newLexer(input string) *lexer {
	l := lexer{
		input: input,
		start: token.Position{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
		tokens: make(chan token.Token, 2),
	}
	l.pos = l.start

	go l.run()
	return &l
}

func TokensFromString(input string) []token.Token {
	l := newLexer(input)

	var tokens []token.Token
	for {
		tok := <-l.tokens
		if tok.Kind == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens
}

func ChanFromString(input string) <-chan token.Token {
	l := newLexer(input)
	return l.tokens
}
