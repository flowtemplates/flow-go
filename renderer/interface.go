package renderer

import (
	"strings"

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

	return render(ast, context)
}

func render(ast []parser.Node, context Context) (string, error) {
	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case parser.Text:
			for _, s := range n.Val {
				result.WriteString(s)
			}
		case parser.ExprBlock:
			s, err := exprToValue(n.Body, context)
			if err != nil {
				return "", err
			}

			result.WriteString(s.String())
		case parser.IfStmt:
			conditionValue, err := exprToValue(n.BegTag.Body, context)
			if err != nil {
				return "", err
			}

			if conditionValue.Boolean() {
				bodyContent, err := render(n.Body, context)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else if n.Else != nil {
				elseContent, err := render(n.Else, context)
				if err != nil {
					return "", err
				}

				result.WriteString(elseContent)
			}
		}
	}

	return result.String(), nil
}

func RenderString(input string, scope Scope) (string, error) {
	tokens := lexer.TokensFromString(input)
	ast, errs := parser.New(tokens).Parse()
	if len(errs) != 0 {
		return "", errs[0]
	}

	res, err := RenderAst(ast, scope)
	if err != nil {
		return "", err
	}

	return res, nil
}
