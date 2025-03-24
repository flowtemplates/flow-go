package parser

import (
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

type Ast []Node

type Node interface {
	node()
}

type Expr interface {
	expr()
}

type Stmt interface {
	Node
	stmt()
}

type (
	NumberLit struct {
		Pos   token.Position
		Value value.NumberValue
	}

	StringLit struct {
		Pos   token.Position
		Quote byte
		Value value.StringValue
	}

	Ident struct {
		Pos  token.Position
		Name string
	}

	Kw struct {
		Kind token.Kind
		Pos  token.Position
	}

	UnaryExpr struct {
		Op   Kw
		Expr Expr
	}

	BinaryExpr struct {
		X  Expr
		Op Kw
		Y  Expr
	}

	TernaryExpr struct {
		Condition Expr
		Do        Kw
		TrueExpr  Expr
		Else      Kw
		FalseExpr Expr
	}

	ParenExpr struct {
		Expr
		Lparen token.Position
		Rparen token.Position
	}

	FilterExpr struct {
		Expr
		OpPos  token.Position
		Filter Ident
	}

	StmtTag struct {
		PreWs string
		// LStmt token.Position
		// RStmt token.Position
	}

	StmtTagWithKw struct {
		StmtTag
		Kw
	}

	StmtTagWithExpr struct {
		StmtTag
		Expr
	}

	ElseIfNode struct {
		ElseIfTag StmtTagWithExpr
		Body      []Node
	}

	ElseNode struct {
		ElseTag StmtTag
		Body    []Node
	}

	CaseClause struct {
		CaseTag StmtTagWithExpr
		Body    []Node
	}
)

// Nodes
type (
	CommNode struct {
		PreWs string
		Pos   token.Position
		// TODO: store val without spaces on the sides to format properly
		Val    string
		PostLB string
	}

	TextNode struct {
		Pos token.Position
		Val []string
	}

	ExprNode struct {
		// LBrace token.Position
		Body Expr
		// RBrace token.Position
	}

	StmtNode struct {
		StmtTagWithKw
		Expr
	}

	IfNode struct {
		IfTag    StmtTagWithExpr
		MainBody []Node
		ElseIfs  []ElseIfNode
		ElseBody ElseNode
		EndTag   StmtTag
	}

	SwitchNode struct {
		SwitchTag   StmtTagWithExpr
		Cases       []CaseClause
		DefaultCase []Node
		EndTag      StmtTag
	}
)

func (*CommNode) node()   {}
func (*TextNode) node()   {}
func (*ExprNode) node()   {}
func (*StmtNode) node()   {}
func (*IfNode) node()     {}
func (*SwitchNode) node() {}

// exprNode() ensures that only expression/type nodes can be
// assigned to an Expr.
func (*NumberLit) expr()   {}
func (*StringLit) expr()   {}
func (*Ident) expr()       {}
func (*UnaryExpr) expr()   {}
func (*BinaryExpr) expr()  {}
func (*TernaryExpr) expr() {}
func (*ParenExpr) expr()   {}
func (*FilterExpr) expr()  {}

// stmtNode() ensures that only statement nodes can be
// assigned to a Stmt.
func (*IfNode) stmt()          {}
func (*StmtTagWithExpr) stmt() {}
func (*SwitchNode) stmt()      {}
