package renderer

import (
	"github.com/flowtemplates/flow-go/lexer"
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

	return render(ast, "", context)
}

func RenderString(input string, scope Scope) (string, error) {
	tokens := lexer.TokensFromString(input)
	ast, err := parser.New(tokens).Parse()
	if err != nil {
		return "", err
	}

	res, err := RenderAst(ast, scope)
	if err != nil {
		return "", err
	}

	return res, nil
}
