package analyzer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
)

type Analyzer struct {
	Tm   TypeMap
	Errs TypeErrors
}

func New() *Analyzer {
	return &Analyzer{
		Tm:   TypeMap{},
		Errs: TypeErrors{},
	}
}

func Typecheck(scope renderer.Input, tm TypeMap) []TypeError {
	errs := []TypeError{}

	for name, typ := range tm {
		prim := tm.getPrimitive(typ)
		if prim == nil {
			continue
		}

		value, ok := scope[name]
		if !ok {
			scope[name] = prim
		} else if !prim.IsValid(value) {
			errs = append(errs, TypeError{
				ExpectedType: *prim,
				Name:         name,
			})
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

// TODO: make func that returns TypeMap and TypeErrors
func (a *Analyzer) TypeMapFromAst(ast []parser.Node) {
	a.parseNodes(ast)
	//	if len(errs) > 0 {
	//		return &errs
	//	}
	//
	// return nil
}

func (a *Analyzer) TypeMapFromBytes(input []byte) error {
	ast, err := parser.AstFromBytes(input)
	if err != nil {
		return fmt.Errorf("ast from bytes: %w", err)
	}

	a.TypeMapFromAst(ast)

	return nil
}
