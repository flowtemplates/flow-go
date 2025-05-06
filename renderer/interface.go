package renderer

import (
	"github.com/flowtemplates/flow-go/filetree"
	"github.com/flowtemplates/flow-go/parser"
)

type File struct {
	Name   string
	Source string
}

type Dir struct {
	Name  string
	Files []File
	Dirs  []Dir
}

func RenderFileTree(ft *filetree.FileTree, input Input) (Dir, error) {
	return Dir{}, nil
}

func RenderTemplate(ast parser.AST, input Input) (string, error) {
	context := InputToContext(input)

	return render(ast, context)
}
