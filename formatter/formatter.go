package formatter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
)

func FromBytes(input []byte) ([]byte, error) {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return nil, fmt.Errorf("ast parsing: %w", err)
	}

	return FromAst(ast)
}

func FromAst(ast parser.Ast) ([]byte, error) {
	var buf bytes.Buffer
	for _, node := range ast {
		if err := formatNode(node, 0, &buf); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func writeSpace(buf *bytes.Buffer) {
	buf.WriteRune(' ')
}

func writeToken(buf *bytes.Buffer, kind token.Kind) {
	buf.WriteString(token.TokenString(kind))
}

func formatNode(node parser.Node, indentLevel int, buf *bytes.Buffer) error {
	indent := strings.Repeat("\t", indentLevel)

	switch n := node.(type) {
	case *parser.TextNode:
		buf.WriteString(indent)
		buf.WriteString(strings.Join(n.Val, ""))
	case *parser.ExprNode:
		writeToken(buf, token.LEXPR)
		writeSpace(buf)

		if err := formatExpr(n.Body, buf); err != nil {
			return err
		}

		writeSpace(buf)
		writeToken(buf, token.REXPR)
	default:
		return fmt.Errorf("unknown node type: %s", n)
	}

	return nil
}

func formatExpr(expr parser.Expr, buf *bytes.Buffer) error {
	switch nn := expr.(type) {
	case *parser.StringLit:
		buf.WriteByte(nn.Quote)
		buf.WriteString(nn.Value.AsString())
		buf.WriteByte(nn.Quote)
	case *parser.NumberLit:
		buf.WriteString(nn.Value.AsString())
	case *parser.ParenExpr:
		writeToken(buf, token.LPAREN)
		if err := formatExpr(nn.Expr, buf); err != nil {
			return err
		}

		writeToken(buf, token.RPAREN)
	case *parser.UnaryExpr:
		buf.WriteString(nn.Op.Kind.String())
		if nn.Op.Kind.IsOneOfMany(token.NOT) {
			writeSpace(buf)
		}

		if err := formatExpr(nn.Expr, buf); err != nil {
			return err
		}
	case *parser.BinaryExpr:
		if err := formatExpr(nn.X, buf); err != nil {
			return err
		}

		writeSpace(buf)
		buf.WriteString(nn.Op.Kind.String())
		writeSpace(buf)

		if err := formatExpr(nn.Y, buf); err != nil {
			return err
		}
	case *parser.Ident:
		buf.WriteString(nn.Name)
	case *parser.TernaryExpr:
		if err := formatExpr(nn.Condition, buf); err != nil {
			return err
		}

		writeSpace(buf)
		buf.WriteString(nn.Do.Kind.String())
		writeSpace(buf)

		if err := formatExpr(nn.TrueExpr, buf); err != nil {
			return err
		}

		writeSpace(buf)
		buf.WriteString(nn.Else.Kind.String())
		writeSpace(buf)

		if err := formatExpr(nn.FalseExpr, buf); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown expression type: %v", nn)
	}

	return nil
}
