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
	f.buf.WriteString(kind.String())
}

func (f *formatter) writeClause(preWs string, tokens ...token.Kind) {
	f.buf.WriteString(preWs)
	f.writeToken(token.LSTMT)
	f.writeSpace()

	for _, t := range tokens {
		f.writeToken(t)
		f.writeSpace()
	}

	f.writeToken(token.RSTMT)
	f.writeLineBreak()
}

func (f *formatter) writeClauseWithExpr(preWs string, expr parser.Expr, tokens ...token.Kind) error {
	f.buf.WriteString(preWs)
	f.writeToken(token.LSTMT)
	f.writeSpace()

	for _, t := range tokens {
		f.writeToken(t)
		f.writeSpace()
	}

	if err := f.writeExpr(expr); err != nil {
		return err
	}

	f.writeSpace()
	f.writeToken(token.RSTMT)
	f.writeLineBreak()

	return nil
}

func (f *formatter) writeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.TextNode:
		f.buf.WriteString(strings.Join(n.Val, ""))

	case *parser.CommNode:
		f.buf.WriteString(n.PreWs)
		f.writeToken(token.LCOMM)
		f.writeSpace()

		f.buf.WriteString(n.Val)

		f.writeSpace()
		f.writeToken(token.RCOMM)

		f.buf.WriteString(n.PostLB)

	case *parser.ExprNode:
		f.writeToken(token.LEXPR)
		f.writeSpace()

		if err := f.writeExpr(n.Body); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(token.REXPR)

	case *parser.GenifNode:
		if err := f.writeClauseWithExpr(n.PreWs, n.Expr, token.GENIF); err != nil {
			return err
		}

	case *parser.IfNode:
		if err := f.writeClauseWithExpr(n.IfTag.PreWs, n.IfTag.Expr, token.IF); err != nil {
			return err
		}

		for _, node := range n.Main {
			if err := f.writeNode(node); err != nil {
				return err
			}
		}

		for _, elseIf := range n.ElseIfs {
			if err := f.writeClauseWithExpr(elseIf.Tag.PreWs, elseIf.Tag.Expr, token.ELSE, token.IF); err != nil {
				return err
			}

			for _, node := range elseIf.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		if len(n.Else.Body) > 0 {
			f.writeClause(n.Else.Tag.PreWs, token.ELSE)

			for _, node := range n.Else.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		f.writeClause(n.EndTag.PreWs, token.END)

	case *parser.SwitchNode:
		if err := f.writeClauseWithExpr(n.SwitchTag.PreWs, n.SwitchTag.Expr, token.SWITCH); err != nil {
			return err
		}

		for _, cc := range n.Cases {
			if err := f.writeClauseWithExpr(cc.Tag.PreWs, cc.Tag.Expr, token.CASE); err != nil {
				return err
			}

			for _, node := range cc.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		if n.DefaultCase != nil {
			f.writeClause(n.DefaultCase.Tag.PreWs, token.DEFAULT)

			for _, node := range n.DefaultCase.Body {
				if err := f.writeNode(node); err != nil {
					return err
				}
			}
		}

		f.writeClause(n.EndTag.PreWs, token.END)

	default:
		return fmt.Errorf("unknown node type: %s", n)
	}

	return nil
}

func (f *formatter) writeExpr(expr parser.Expr) error {
	switch e := expr.(type) {
	case *parser.StringLit:
		f.buf.WriteByte(e.Quote)
		f.buf.WriteString(e.Value.AsString())
		f.buf.WriteByte(e.Quote)

	case *parser.NumberLit:
		f.buf.WriteString(e.Value.AsString())

	case *parser.ParenExpr:
		f.writeToken(token.LPAREN)

		if err := f.writeExpr(e.Expr); err != nil {
			return err
		}

		f.writeToken(token.RPAREN)

	case *parser.UnaryExpr:
		f.writeToken(e.Op.Kind)

		if e.Op.Kind.IsOneOfMany(token.NOT) {
			f.writeSpace()
		}

		if err := f.writeExpr(e.Expr); err != nil {
			return err
		}

	case *parser.BinaryExpr:
		if err := f.writeExpr(e.X); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(e.Op.Kind)
		f.writeSpace()

		if err := f.writeExpr(e.Y); err != nil {
			return err
		}

	case *parser.Ident:
		f.buf.WriteString(e.Name)

	case *parser.TernaryExpr:
		if err := f.writeExpr(e.Condition); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(e.Do.Kind)
		f.writeSpace()

		if err := f.writeExpr(e.TrueExpr); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(e.Else.Kind)
		f.writeSpace()

		if err := f.writeExpr(e.FalseExpr); err != nil {
			return err
		}

	case *parser.FilterExpr:
		if err := f.writeExpr(e.Expr); err != nil {
			return err
		}

		f.writeSpace()
		f.writeToken(token.RARR)
		f.writeSpace()

		f.buf.WriteString(e.Filter.Name)

	default:
		return fmt.Errorf("unknown expression type: %T", e)
	}

	return nil
}
