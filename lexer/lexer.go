package lexer

import (
	"github.com/flowtemplates/flow-go/token"
)

const eof = 0

func (l *lexer) emit(t token.Kind) {
	if l.startPos.Offset < l.pos.Offset {
		l.tokensCh <- token.Token{
			Kind: t,
			Val:  string(l.source[l.startPos.Offset:l.pos.Offset]),
			Pos:  l.startPos,
		}
		l.startPos = l.pos
	}
}

func (l *lexer) run() {
	for state := lexLineWhitespace(lexText); state != nil; {
		state = state(l)
	}
	close(l.tokensCh)
}

func (l *lexer) next() rune {
	if l.pos.Offset >= len(l.source) {
		return eof
	}

	r := rune(l.source[l.pos.Offset])
	l.pos.Offset++

	// Handle newline properly
	if r == '\n' {
		l.pos.Line++
		l.pos.Column = 1 // Reset column on new line
	} else {
		l.pos.Column++
	}

	return r
}

func (l *lexer) back() {
	if l.pos.Offset <= 0 {
		return
	}

	l.pos.Offset--
	if l.source[l.pos.Offset] == '\n' {
		l.pos.Line-- // Move back a line
		// Find previous line start to restore column correctly
		l.pos.Column = 1
		for i := l.pos.Offset - 1; i >= 0 && l.source[i] != '\n'; i-- {
			l.pos.Column++
		}
	} else {
		l.pos.Column--
	}
}

func (l *lexer) peek() rune {
	if l.pos.Offset < len(l.source) {
		r := rune(l.source[l.pos.Offset])
		return r
	}
	return eof
}

func (l *lexer) accept(valid string) bool {
	r := l.next()
	for _, c := range valid {
		if c == r {
			return true
		}
	}

	l.back()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}
