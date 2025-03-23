package parser

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/token"
)

type ErrorType string

const (
	ErrExpressionExpected ErrorType = "expression expected"
	// TODO: change message
	// ErrUnexpectedBeforeStmt ErrorType = "unexpected text before statement tag"
	ErrEndExpected     ErrorType = "'{% end %}' expected"
	ErrKeywordExpected ErrorType = "'if', 'genif', 'end' expected"
)

type Error struct {
	Pos token.Position
	Typ ErrorType
}

func (e Error) Error() string {
	return string(e.Typ)
}

type ExpectedTokensError struct {
	Pos    token.Position
	Tokens []token.Kind
}

func (e ExpectedTokensError) Error() string {
	b := []string{}
	for _, e := range e.Tokens {
		b = append(b, fmt.Sprintf("'%s'", e.String()))
	}

	return fmt.Sprintf("%s expected", strings.Join(b, ", "))
}

// type ErrorList []error // nolint: errname

// func (l ErrorList) Error() string {
// 	switch len(l) {
// 	case 0:
// 		return "no errors"
// 	case 1:
// 		return l[0].Error()
// 	}

// 	b := []string{}
// 	for _, e := range l {
// 		b = append(b, e.Error())
// 	}

// 	return strings.Join(b, ", ")
// }

// // Err returns an error equivalent to this error list.
// // If the list is empty, Err returns nil.
// func (l ErrorList) Err() error {
// 	if len(l) == 0 {
// 		return nil
// 	}

// 	return l
// }

// // Add adds an [Error] with given position and error message to an [ErrorList].
// func (l *ErrorList) Add(err error) {
// 	*l = append(*l, err)
// }
