package renderer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
)

func RenderAst(ast []parser.Node, scope Scope) (string, error) {
	// tm := make(analyzer.TypeMap)

	// if errs := analyzer.GetTypeMapFromAst(ast, tm); len(errs) != 0 {
	// 	return "", errs[0] // TODO: error handling
	// }

	// if errs := analyzer.Typecheck(scope, tm); len(errs) != 0 {
	// 	return "", errs[0] // TODO: error handling
	// }

	context := scopeToContext(scope)

	return render(ast, context)
}

func RenderBytes(input []byte, scope Scope) (string, error) {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return "", fmt.Errorf("ast from string: %w", err)
	}

	res, err := RenderAst(ast, scope)
	if err != nil {
		return "", err
	}

	return res, nil
}
