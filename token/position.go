package token

import (
	"fmt"
)

type Pos struct {
	Row    int `json:"-"` // FIXME:
	Column int `json:"-"`
	Offset int `json:"-"`
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d:%d", p.Row, p.Column, p.Offset)
}
