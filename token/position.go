package token

import (
	"fmt"
)

// TODO: write equal function for ast
type Position struct {
	Line   int `json:"-"`
	Column int `json:"-"`
	Offset int `json:"-"`
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d:%d", p.Line, p.Column, p.Offset)
}
