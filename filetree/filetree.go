package filetree

import "github.com/flowtemplates/flow-go/parser"

type FileTree Dir

type File struct {
	Path      string
	Name      parser.AST
	Condition parser.Expr
	Content   parser.AST
}

type Dir struct {
	Path      string
	Name      parser.AST
	Condition parser.Expr
	Files     []File
	Dirs      []Dir
}
