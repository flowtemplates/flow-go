package analyzer

import (
	"github.com/flowtemplates/flow-go/filetree"
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
)

type analyzer struct {
	Tm   TypeMap
	Errs TypeErrors
}

func newAnalyzer() *analyzer {
	return &analyzer{
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

func TypeMapFromAST(ast parser.AST) (TypeMap, TypeErrors) {
	a := newAnalyzer()

	a.parseNodes(ast)

	if len(a.Errs) == 0 {
		return a.Tm, nil
	}

	return a.Tm, a.Errs
}

func TypeMapFromFileTree(ft *filetree.FileTree) (TypeMap, TypeErrors) {
	a := newAnalyzer()

	a.analyzeDir((*filetree.Dir)(ft))

	if len(a.Errs) == 0 {
		return a.Tm, nil
	}

	return a.Tm, a.Errs
}

func (a analyzer) analyzeDir(d *filetree.Dir) {
	for _, dir := range d.Dirs {
		a.analyzeDir(&dir)
	}

	for _, file := range d.Files {
		a.parseNodes(file.Content)
	}
}
