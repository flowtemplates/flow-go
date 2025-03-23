package lexer

import (
	"bytes"
	"unicode"

	"github.com/flowtemplates/flow-go/token"
)

type stateFn func(*lexer) stateFn

func (l *lexer) lexToken(t token.Kind, next stateFn) stateFn {
	tokLen := len(token.TokenString(t))
	l.pos.Offset += tokLen
	l.pos.Column += tokLen
	l.emit(t)
	return next
}

func (l *lexer) startsWith(t token.Kind) bool {
	tokBytes := token.TokenBytes(t)
	if len(tokBytes) > 0 {
		return bytes.HasPrefix(l.source[l.pos.Offset:], tokBytes)
	}

	return false
}

func (l *lexer) tryTokens(nextState stateFn, tokens ...token.Kind) stateFn {
	for _, t := range tokens {
		if l.startsWith(t) {
			return l.lexToken(t, nextState)
		}
	}

	return nil
}

// TODO: rewrite this whole CRAP
func (l *lexer) tryKeywords(nextState stateFn) stateFn {
	a := l.source[l.pos.Offset:]
	for _, tok := range token.GetKeywords() {
		tokBytes := token.TokenBytes(tok)
		if len(tokBytes) > 0 && bytes.HasPrefix(a, tokBytes) {
			if len(l.source) < l.pos.Offset+len(tokBytes)+1 {
				return l.lexToken(tok, nextState)
			}

			if unicode.IsSpace(rune(l.source[l.pos.Offset+len(tokBytes)])) {
				return l.lexToken(tok, nextState)
			}

			b := l.source[l.pos.Offset+len(tokBytes):]
			for _, tok2 := range token.GetOperators() {
				tokBytes2 := token.TokenBytes(tok2)
				if len(tokBytes2) > 0 && bytes.HasPrefix(b, tokBytes2) {
					return l.lexToken(tok, nextState)
				}
			}
		}
	}

	return nil
}

func lexText(l *lexer) stateFn {
	for {
		r := l.peek()
		if r == eof {
			l.emit(token.TEXT)
			return nil
		}

		if r == '\n' {
			l.emit(token.TEXT)
			l.next()
			l.emit(token.LNBR)
			return lexLineWhitespace(lexText)
		}

		// if unicode.IsSpace(r) {
		// 	return lexLineWhitespace(lexText)
		// }

		if l.startsWith(token.LEXPR) {
			l.emit(token.TEXT)
			return l.lexToken(token.LEXPR, lexExpr)
		}

		if l.startsWith(token.RARR) {
			l.emit(token.TEXT)
			return l.lexToken(token.RARR, lexComm)
		}

		if l.startsWith(token.LSTMT) {
			l.emit(token.TEXT)
			return l.lexToken(token.LSTMT, lexStmt)
		}

		if l.startsWith(token.LCOMM) {
			l.emit(token.TEXT)
			return l.lexToken(token.LCOMM, lexComm)
		}

		if l.startsWith(token.REXPR) {
			l.emit(token.TEXT)
			return l.lexToken(token.REXPR, lexText)
		}

		if l.startsWith(token.RCOMM) {
			l.emit(token.TEXT)
			return l.lexToken(token.RCOMM, lexLineWhitespace(lexText))
		}

		if l.startsWith(token.RSTMT) {
			l.emit(token.TEXT)
			return l.lexToken(token.RSTMT, lexLineWhitespace(lexText))
		}

		l.next()
	}
}

// TODO: rename
func lexRealExpr(nextState stateFn) stateFn {
	return func(l *lexer) stateFn {
		r := l.next()

		if r == eof {
			return nil
		}
		if r == '\n' || r == '\r' {
			l.back()
			return lexText
		}
		if unicode.IsSpace(r) {
			return lexLineWhitespace(nextState)
		}
		if r == '\'' {
			return lexSQString
		}
		if r == '"' {
			return lexDQString
		}
		if r == token.TokenRune(token.LPAREN) {
			l.back()
			return l.lexToken(token.LPAREN, nextState)
		}
		if r == token.TokenRune(token.RPAREN) {
			l.back()
			return l.lexToken(token.RPAREN, nextState)
		}
		if unicode.IsDigit(r) {
			return lexNum(nextState)
		}
		if token.IsNotOp(r) && r != token.TokenRune(token.PERIOD) {
			return lexIdent(nextState)
		}

		return nextState
	}
}

func lexExpr(l *lexer) stateFn {
	if l.startsWith(token.REXPR) {
		return l.lexToken(token.REXPR, lexText)
	}

	if state := l.tryTokens(lexExpr, token.GetOperators()...); state != nil {
		return state
	}

	if state := l.tryKeywords(lexExpr); state != nil {
		return state
	}

	return lexRealExpr(lexExpr)
}

func lexComm(l *lexer) stateFn {
	// ? try to lex something to do not cause commenting whole thing if there is no closing tag
	for {
		if l.startsWith(token.RCOMM) {
			l.emit(token.COMM_TEXT)
			return l.lexToken(token.RCOMM, lexLineWhitespace(lexText))
		}

		r := l.next()
		if r == eof {
			l.emit(token.COMM_TEXT)
			return nil
		}
	}
}

func lexNum(nextState stateFn) stateFn {
	return func(l *lexer) stateFn {
		digits := "0123456789"

		l.acceptRun(digits)
		if l.accept(".") {
			l.acceptRun(digits)
			l.emit(token.FLOAT)
		} else {
			l.emit(token.INT)
		}

		return nextState
	}
}

func lexSQString(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case '\n':
			l.back()
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case token.SQUOTE:
			l.emit(token.STR)
			return lexExpr
		}
	}
}

// TODO: refactor
func lexDQString(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case '\n':
			l.back()
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case token.DQUOTE:
			l.emit(token.STR)
			return lexExpr
		}
	}
}

func lexIdent(nextState stateFn) stateFn {
	return func(l *lexer) stateFn {
		for {
			switch r := l.next(); {
			case r == eof:
				l.emit(token.IDENT)
				return nil
			case !token.IsNotOp(r) || unicode.IsSpace(r):
				l.back()
				l.emit(token.IDENT)
				return nextState
			}
		}
	}
}

func lexLineWhitespace(nextState stateFn) stateFn {
	return func(l *lexer) stateFn {
		for {
			switch r := l.peek(); {
			case r == ' ' || r == '\t':
				l.next()
			case unicode.IsSpace(r):
				l.emit(token.WS)
				return lexText
			default:
				l.emit(token.WS)
				return nextState
			}
		}
	}
}

func lexStmt(l *lexer) stateFn {
	if l.startsWith(token.RSTMT) {
		return l.lexToken(token.RSTMT, lexLineWhitespace(lexText))
	}

	if state := l.tryTokens(lexStmt, token.GetOperators()...); state != nil {
		return state
	}

	if state := l.tryKeywords(lexStmt); state != nil {
		return state
	}

	return lexRealExpr(lexStmt)
}
