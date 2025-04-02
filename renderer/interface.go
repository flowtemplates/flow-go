package renderer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
)

func RenderAst(ast []parser.Node, scope Scope) ([]byte, error) {
	// tm := make(analyzer.TypeMap)

	// if errs := analyzer.GetTypeMapFromAst(ast, tm); len(errs) != 0 {
	// 	return "", errs[0] // TODO: error handling
	// }

	// if errs := analyzer.Typecheck(scope, tm); len(errs) != 0 {
	// 	return "", errs[0] // TODO: error handling
	// }

	context := ScopeToContext(scope)

	return render(ast, context)
}

func RenderBytes(input []byte, scope Scope) ([]byte, error) {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return nil, fmt.Errorf("ast from bytes: %w", err)
	}

	res, err := RenderAst(ast, scope)
	if err != nil {
		return nil, err
	}

	return res, nil
}
