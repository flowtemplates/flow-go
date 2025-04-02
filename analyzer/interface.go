package analyzer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
)

// func Typecheck(scope renderer.Scope, tm TypeMap) []TypeError {
// 	errs := []TypeError{}
// 	for name, typ := range tm {
// 		if typ == types.Any {
// 			continue
// 		}

// 		value, ok := scope[name]
// 		if !ok {
// 			scope[name] = typ.GetDefaultValue()
// 		} else if !typ.IsValid(value) {
// 			errs = append(errs, TypeError{
// 				ExpectedType: typ,
// 				Name:         name,
// 				Val:          value,
// 			})
// 		}
// 	}

// 	if len(errs) != 0 {
// 		return errs
// 	}

// 	return nil
// }

func TypeMapFromAst(ast []parser.Node, tm TypeMap, scope renderer.Scope) *TypeErrors {
	context := renderer.ScopeToContext(scope)
	errs := TypeErrors{}
	parseNodes(ast, tm, context, &errs)

	if len(errs) > 0 {
		return &errs
	}

	return nil
}

func TypeMapFromBytes(input []byte, tm TypeMap, scope renderer.Scope) error {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return fmt.Errorf("ast from bytes: %w", err)
	}

	if errs := TypeMapFromAst(ast, tm, scope); errs != nil {
		return errs
	}

	return nil
}
