package lexer

import (
	"strings"
	"unicode"

	"github.com/flowtemplates/flow-go/token"
)

type stateFn func(*Lexer) stateFn

func (l *Lexer) lexToken(t token.Kind, next stateFn) stateFn {
	tokLen := len(token.TokenString(t))
	l.pos.Offset += tokLen
	l.pos.Column += tokLen
	l.emit(t)
	return next
}

func (l *Lexer) StartsWith(t token.Kind) bool {
	tokString := token.TokenString(t)
	if tokString != "" {
		return strings.HasPrefix(l.input[l.pos.Offset:], tokString)
	}

	return false
}

func (l *Lexer) tryTokens(nextState stateFn, tokens ...token.Kind) stateFn {
	for _, token := range tokens {
		if l.StartsWith(token) {
			return l.lexToken(token, nextState)
		}
	}

	return nil
}

func lexText(l *Lexer) stateFn {
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

		if l.StartsWith(token.LEXPR) {
			l.emit(token.TEXT)
			return l.lexToken(token.LEXPR, lexExpr)
		}

		if l.StartsWith(token.RARR) {
			l.emit(token.TEXT)
			return l.lexToken(token.RARR, lexComm)
		}

		if l.StartsWith(token.LSTMT) {
			l.emit(token.TEXT)
			return l.lexToken(token.LSTMT, lexStmt)
		}

		if l.StartsWith(token.LCOMM) {
			l.emit(token.TEXT)
			return l.lexToken(token.LCOMM, lexComm)
		}

		if l.StartsWith(token.REXPR) {
			l.emit(token.TEXT)
			return l.lexToken(token.REXPR, lexLineWhitespace(lexText))
		}

		if l.StartsWith(token.RCOMM) {
			l.emit(token.TEXT)
			return l.lexToken(token.RCOMM, lexText)
		}

		if l.StartsWith(token.RSTMT) {
			l.emit(token.TEXT)
			return l.lexToken(token.RSTMT, lexLineWhitespace(lexText))
		}

		l.next()
	}
}

// TODO: rename
func lexRealExpr(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
		switch r := l.next(); {
		case r == eof:
			return nil
		case r == '\n' || r == '\r':
			l.back()
			return lexText
		case unicode.IsSpace(r):
			return lexLineWhitespace(nextState)
		case r == '"' || r == '\'':
			return lexString
		case r == token.TokenRune(token.LPAREN):
			l.back()
			return l.lexToken(token.LPAREN, nextState)
		case r == token.TokenRune(token.RPAREN):
			l.back()
			return l.lexToken(token.RPAREN, nextState)
		case unicode.IsDigit(r):
			return lexNum(nextState)
		case token.IsNotOp(r) && r != '.':
			return lexIdent(nextState)
		default:
			l.emit(token.EXPECTED_EXPR)
			return nextState
		}
	}
}

func lexExpr(l *Lexer) stateFn {
	if l.StartsWith(token.REXPR) {
		return l.lexToken(token.REXPR, lexText)
	}

	if state := l.tryTokens(lexExpr, append(token.GetOperators(),
		token.NOT,
		token.AND,
		token.DO,
		token.ELSE,
		token.OR,
		token.IS,
	)...); state != nil {
		return state
	}

	return lexRealExpr(lexExpr)
}

func lexComm(l *Lexer) stateFn {
	// ? try to lex something to do not cause commenting whole thing if there is no closing tag
	for {
		if l.StartsWith(token.RCOMM) {
			l.emit(token.COMM_TEXT)
			return l.lexToken(token.RCOMM, lexText)
		}

		r := l.next()
		if r == eof {
			l.emit(token.COMM_TEXT)
			return nil
		}
	}
}

func lexNum(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
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
func lexString(l *Lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			l.emit(token.NOT_TERMINATED_STR)
			return lexText
		case '"', '\'':
			l.emit(token.STR)
			return lexExpr
		}
	}
}

func lexIdent(nextState stateFn) stateFn {
	return func(l *Lexer) stateFn {
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
	return func(l *Lexer) stateFn {
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

func lexStmt(l *Lexer) stateFn {
	if l.StartsWith(token.RSTMT) {
		return l.lexToken(token.RSTMT, lexLineWhitespace(lexText))
	}

	if state := l.tryTokens(lexStmt, append(token.GetOperators(), token.GetKeywords()...)...); state != nil {
		return state
	}

	return lexRealExpr(lexStmt)
}
