package renderer

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/lexer"
	"github.com/flowtemplates/flow-go/parser"
)

func RenderAst(ast []parser.Node, scope Scope) (string, error) {
	tm := make(analyzer.TypeMap)

	if errs := analyzer.GetTypeMapFromAst(ast, tm); len(errs) != 0 {
		return "", errs[0] // TODO: error handling
	}

	if errs := analyzer.Typecheck(scope, tm); len(errs) != 0 {
		return "", errs[0] // TODO: error handling
	}

	context := scopeToContext(scope, tm)

	var result strings.Builder
	for _, node := range ast {
		switch n := node.(type) {
		case parser.Text:
			for _, s := range n.Val {
				result.WriteString(s)
			}
		case parser.ExprBlock:
			switch body := n.Body.(type) {
			case parser.Ident:
				value, exists := context[body.Name]
				if !exists {
					return "", fmt.Errorf("%s not declared", body.Name)
				}

				result.WriteString(value.String())
			case parser.Lit:
				result.WriteString(body.Val)
			default:
				return "", fmt.Errorf("unsupported expr type: %T", body)
			}
		case parser.IfStmt:
			conditionValue, err := exprToValue(n.BegTag.Body, context)
			if err != nil {
				return "", err
			}

			if conditionValue.Boolean() {
				bodyContent, err := RenderAst(n.Body, scope)
				if err != nil {
					return "", err
				}

				result.WriteString(bodyContent)
			} else if n.Else != nil {
				elseContent, err := RenderAst(n.Else, scope)
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
		return "", fmt.Errorf("failed to render: %w", err)
	}

	return res, nil
}
