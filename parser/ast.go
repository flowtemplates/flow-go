package parser

import (
	"github.com/flowtemplates/flow-go/token"
)

type AST []Node

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	// End() token.Pos // position of first character immediately after the node
}

type Expr interface {
	Node
	expr()
}

type (
	// A BadExpr node is a placeholder for an expression containing
	// syntax errors for which a correct expression node cannot be
	// created.
	BadExpr struct {
		From, To token.Pos
	}

	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT, token.FLOAT, token.BOOL, token.STRING or token.NULL
		Value    string      // literal string; e.g. 42, 3.14, 'a', "foo", true, false
	}

	// TODO:
	// Array struct {
	// 	Elems  []Expr
	// 	Lbrack token.Pos // position of "["
	// 	Rbrack token.Pos // position of "]"
	// }
	//
	// Object struct {
	// 	Pairs  map[string]Expr
	// 	Lbrack token.Pos // position of "{"
	// 	Rbrack token.Pos // position of "}"
	// }
	//
	// // Indexing: arr[0]
	// IndexExpr struct {
	// 	Target Expr
	// 	Index  Expr
	// 	Pos_   int
	// }
	// // Field access: obj.a
	// FieldExpr struct {
	// 	Target Expr
	// 	Field  string
	// 	Pos_   int
	// }

	Ident struct {
		Name    string
		NamePos token.Pos
	}

	UnaryExpr struct {
		Op    token.Kind
		OpPos token.Pos
		Expr  Expr
	}

	BinaryExpr struct {
		X     Expr
		Op    token.Kind
		OpPos token.Pos
		Y     Expr
	}

	TernaryExpr struct {
		Condition Expr
		ThenOp    token.Kind
		ThenExpr  Expr
		ElseOp    token.Kind
		ElseExpr  Expr
	}

	ParenExpr struct {
		Lparen token.Pos
		X      Expr
		Rparen token.Pos
	}

	PipeExpr struct {
		X       Expr
		Filters []Filter
	}

	TypeAssertExpr struct {
		Type    string
		TypePos token.Pos
		X       Expr
	}
)

type Filter struct {
	Name string
	Args []Expr
}

type Stmt interface {
	Node
	stmt()
}

type (
	Comment struct {
		Lbrace token.Pos
		Rbrace token.Pos
		Text   string
	}

	Text struct {
		Value string
		Pos_  token.Pos
	}

	Print struct {
		Expr   Expr
		Lbrace token.Pos
		Rbrace token.Pos
	}

	If struct {
		Conditions []IfBranch
		ElseBody   AST
	}

	IfBranch struct {
		Condition Expr
		Body      AST
		Lbrace    token.Pos
		Rbrace    token.Pos
	}

	// ForStmt struct {
	// 	KeyVar, ValueVar string
	// 	Iterable         Expr
	// 	Block            []Stmt
	// 	Pos_             int
	// }

	Switch struct {
		Expr    Expr
		Cases   []CaseClause
		Default AST
		Pos_    token.Pos
	}

	CaseClause struct {
		Match Expr
		Body  AST
		Pos_  token.Pos
	}

	Let struct {
		Lhs  Expr
		Rhs  Expr
		Pos_ token.Pos
	}

	// TODO: macros
	// MacroDecl struct {
	// 	Name       string
	// 	Parameters []MacroParam
	// 	Body       []Stmt
	// 	Pos_       int
	// }
	// MacroCall struct {
	// 	Name string
	// 	Args []Expr
	// 	Pos_ int
	// }

	// TODO: template inheritance
	// Extends struct {
	// 	TemplatePath string
	// 	Pos_         int
	// }
	// BlockStmt struct {
	// 	Name string
	// 	Body []Stmt
	// 	Pos_ int
	// }

	Raw struct {
		Content string
		Pos_    token.Pos
	}

	FilterBlock struct {
		FilterName string
		Body       AST
		Pos_       token.Pos
	}
)

// type MacroParam struct {
// 	Name    string
// 	Type    string // optional, for future type checking
// 	Default Expr   // can be nil
// }

func (x *BadExpr) Pos() token.Pos        { return x.From }
func (x *BasicLit) Pos() token.Pos       { return x.ValuePos }
func (x *Ident) Pos() token.Pos          { return x.NamePos }
func (x *UnaryExpr) Pos() token.Pos      { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *TernaryExpr) Pos() token.Pos    { return x.Condition.Pos() }
func (x *ParenExpr) Pos() token.Pos      { return x.Lparen }
func (x *PipeExpr) Pos() token.Pos       { return x.X.Pos() }
func (x *TypeAssertExpr) Pos() token.Pos { return x.TypePos }
func (x *Comment) Pos() token.Pos        { return x.Lbrace }
func (x *Text) Pos() token.Pos           { return x.Pos_ }
func (x *Print) Pos() token.Pos          { return x.Lbrace }
func (x *If) Pos() token.Pos             { return x.Conditions[0].Lbrace }
func (x *IfBranch) Pos() token.Pos       { return x.Lbrace }
func (x *Switch) Pos() token.Pos         { return x.Pos_ }
func (x *CaseClause) Pos() token.Pos     { return x.Pos_ }
func (x *Let) Pos() token.Pos            { return x.Pos_ }
func (x *Raw) Pos() token.Pos            { return x.Pos_ }
func (x *FilterBlock) Pos() token.Pos    { return x.Pos_ }

func (x *BadExpr) expr()        {}
func (x *BasicLit) expr()       {}
func (x *Ident) expr()          {}
func (x *UnaryExpr) expr()      {}
func (x *BinaryExpr) expr()     {}
func (x *TernaryExpr) expr()    {}
func (x *ParenExpr) expr()      {}
func (x *PipeExpr) expr()       {}
func (x *TypeAssertExpr) expr() {}
