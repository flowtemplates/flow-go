package parser

import (
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

type Node interface{} // nolint: iface

type Expr interface{} // nolint: iface

type Stmt interface{} // nolint: iface

type (
	Text struct {
		Pos token.Position
		Val []string
	}

	Lit struct {
		Pos   token.Position
		Value value.Valueable
	}

	Ident struct {
		Pos  token.Position
		Name string
	}

	UnaryExpr struct {
		Op    token.Kind
		OpPos token.Position
		X     Expr
	}

	BinaryExpr struct {
		X     Expr
		Op    token.Kind
		OpPos token.Position
		Y     Expr
	}

	TernaryExpr struct {
		Condition Expr
		Do        token.Kind
		DoPos     token.Position
		TrueExpr  Expr
		Else      token.Kind
		ElsePos   token.Position
		FalseExpr Expr
	}

	StmtTag struct {
		PreWs string
		LStmt token.Position
		Kw    token.Kind
		KwPos token.Position
		RStmt token.Position
	}

	StmtTagWithExpr struct {
		StmtTag
		Body Expr
	}

	ExprBlock struct {
		LBrace token.Position
		Body   Expr
		RBrace token.Position
	}

	IfStmt struct {
		BegTag      StmtTagWithExpr
		Body        []Node
		Else        []Node
		PreEndTagWs string
	}
)
