package formatter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
)

type formatter struct {
	buf *bytes.Buffer
}

func newFormatter() *formatter {
	return &formatter{
		buf: &bytes.Buffer{},
	}
}

func (f *formatter) writeSpace() {
	f.buf.WriteRune(' ')
}

func (f *formatter) writeLineBreak() {
	f.buf.WriteRune('\n')
}

func (f *formatter) writeToken(kind token.Kind) {
	f.buf.WriteString(token.TokenString(kind))
}

func (f *formatter) writeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.TextNode:
		f.buf.WriteString(strings.Join(n.Val, ""))
	case *parser.ExprNode:
		f.writeToken(token.LEXPR)
		f.writeSpace()

		if err := f.writeExpr(n.Body); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(token.REXPR)
	// TODO: refactor
	case *parser.IfNode:
		f.writeToken(token.LSTMT)
		f.writeSpace()

		f.writeToken(token.IF)
		f.writeSpace()

		if err := f.writeExpr(n.IfTag.Expr); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(token.RSTMT)
		f.writeLineBreak()

		for _, node := range n.MainBody {
			if err := f.writeNode(node); err != nil {
				return err
			}
		}

		for _, elseIf := range n.ElseIfs {
			f.writeToken(token.LSTMT)
			f.writeSpace()

			f.writeToken(token.ELSE)
			f.writeSpace()

			f.writeToken(token.IF)
			f.writeSpace()

			if err := f.writeExpr(elseIf.ElseIfTag.Expr); err != nil {
				return err
			}

			f.writeSpace()
			f.writeToken(token.RSTMT)
			f.writeLineBreak()
			for _, node := range elseIf.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		if len(n.ElseBody.Body) > 0 {
			f.writeToken(token.LSTMT)
			f.writeSpace()

			f.writeToken(token.ELSE)
			f.writeSpace()

			f.writeToken(token.RSTMT)
			f.writeLineBreak()

			for _, node := range n.ElseBody.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		f.writeToken(token.LSTMT)
		f.writeSpace()

		f.writeToken(token.END)
		f.writeSpace()

		f.writeToken(token.RSTMT)
		f.writeLineBreak()
	default:
		return fmt.Errorf("unknown node type: %s", n)
	}

	return nil
}

func (f *formatter) writeExpr(expr parser.Expr) error {
	switch nn := expr.(type) {
	case *parser.StringLit:
		f.buf.WriteByte(nn.Quote)
		f.buf.WriteString(nn.Value.AsString())
		f.buf.WriteByte(nn.Quote)
	case *parser.NumberLit:
		f.buf.WriteString(nn.Value.AsString())
	case *parser.ParenExpr:
		f.writeToken(token.LPAREN)
		if err := f.writeExpr(nn.Expr); err != nil {
			return err
		}

		f.writeToken(token.RPAREN)
	case *parser.UnaryExpr:
		f.buf.WriteString(nn.Op.Kind.String())
		if nn.Op.Kind.IsOneOfMany(token.NOT) {
			f.writeSpace()
		}

		if err := f.writeExpr(nn.Expr); err != nil {
			return err
		}
	case *parser.BinaryExpr:
		if err := f.writeExpr(nn.X); err != nil {
			return err
		}

		f.writeSpace()
		f.buf.WriteString(nn.Op.Kind.String())
		f.writeSpace()

		if err := f.writeExpr(nn.Y); err != nil {
			return err
		}
	case *parser.Ident:
		f.buf.WriteString(nn.Name)
	case *parser.TernaryExpr:
		if err := f.writeExpr(nn.Condition); err != nil {
			return err
		}

		f.writeSpace()
		f.buf.WriteString(nn.Do.Kind.String())
		f.writeSpace()

		if err := f.writeExpr(nn.TrueExpr); err != nil {
			return err
		}

		f.writeSpace()
		f.buf.WriteString(nn.Else.Kind.String())
		f.writeSpace()

		if err := f.writeExpr(nn.FalseExpr); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown expression type: %v", nn)
	}

	return nil
}
