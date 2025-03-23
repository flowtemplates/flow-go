package formatter

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
)

func Bytes(input []byte) ([]byte, error) {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return nil, fmt.Errorf("ast parsing: %w", err)
	}

	return Ast(ast)
}

func Ast(ast parser.Ast) ([]byte, error) {
	f := newFormatter()
	for _, node := range ast {
		if err := f.writeNode(node); err != nil {
			return nil, err
		}
	}

	return f.buf.Bytes(), nil
}
