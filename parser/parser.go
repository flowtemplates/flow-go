package parser

import (
	"fmt"
	"strings"

	"github.com/flowtemplates/flow-go/token"
)

type Parser struct {
	tokens  []token.Token
	pos     int
	nodes   []Node
	errors  []error
	current token.Token
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,
	}

	p.next()
	return p
}

func (p *Parser) Parse() ([]Node, []error) {
	for p.pos < len(p.tokens) {
		node := p.parseNode()
		if node != nil {
			p.nodes = append(p.nodes, node)
		} else {
			p.next()
		}
	}

	return p.nodes, p.errors
}

func (p *Parser) errorf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	p.errors = append(p.errors, err)
}

func (p *Parser) next() token.Token {
	p.pos++
	p.current = p.getCurrent()
	return p.current
}

func (p *Parser) back() {
	p.pos--
	p.current = p.getCurrent()
}

func (p *Parser) peek() token.Token {
	t := p.next()
	p.back()
	return t
}

func (p *Parser) consumeWhitespaces() string {
	var builder strings.Builder
	for p.current.Kind == token.WS {
		builder.WriteString(p.current.Val)
		p.next()
	}
	return builder.String()
}

func (p *Parser) consumeLineBreak() string {
	var builder strings.Builder
	for p.current.Kind == token.LNBR {
		builder.WriteString(p.current.Val)
		p.next()
	}
	return builder.String()
}

func (p *Parser) getCurrent() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Kind: token.EOF}
}

func (p *Parser) parseNode() Node {
	switch p.current.Kind {
	case token.TEXT:
		return p.parseText()
	case token.LNBR:
		return p.parseText()
	case token.WS:
		return p.parseText()
	case token.LEXPR:
		return p.parseExprBlock()
	case token.LSTMT:
		return p.parseStmt("")
	// case token.LCOMM:
	// 	return p.parseco()
	case token.EOF:
		return nil // End of input
	default:
		p.errorf("unexpected token: %v", p.current)
		return nil
	}
}

func (p *Parser) parseText() Node {
	var res []string
loop:
	for {
		if p.current.Kind == token.EOF {
			break
		}

		switch p.current.Kind {
		case token.TEXT:
			res = append(res, p.current.Val)
		case token.LNBR:
			res = append(res, p.current.Val)
		case token.WS:
			// switch p.peek().Typ {
			// case token.LSTMT:
			// 	p.next()
			// 	return p.parseStmt(p.current.Val)
			// }
			res = append(res, p.current.Val)
		default:
			break loop
		}

		p.next()
	}

	text := Text{
		Pos: p.current.Pos,
		Val: res,
	}
	return text
}

func (p *Parser) parseExprBlock() Node {
	exprBlock := ExprBlock{
		LBrace: p.current.Pos,
	}
	p.next() // Consume LEXPR
	p.consumeWhitespaces()

	exprBlock.Body = p.parseExpr()

	if p.current.Kind != token.REXPR {
		p.errorf("expected REXPR, got %v", p.current)
		return exprBlock
	}

	exprBlock.RBrace = p.current.Pos
	p.next() // Consume REXPR
	return exprBlock
}

func (p *Parser) parseExpr() Node {
	return p.parseTernaryExpr(1)
}

// parseTernaryExpr parses expressions with ternary operators and precedence.
func (p *Parser) parseTernaryExpr(minPrecedence int) Node {
	condition := p.parseBinaryExpr(minPrecedence)

	doTok := p.current
	if doTok.IsOneOfMany(token.QUESTION, token.DO) {
		ternary := &TernaryExpr{
			Condition: condition,
			DoPos:     p.current.Pos,
			Do:        doTok.Kind,
		}
		p.next() // Consume '?'
		p.consumeWhitespaces()

		ternary.TrueExpr = p.parseExpr()

		var expectedElseTok token.Kind
		if doTok.Kind == token.QUESTION {
			expectedElseTok = token.COLON
		} else {
			expectedElseTok = token.ELSE
		}

		if p.current.Kind != expectedElseTok {
			p.errorf("expected %s in ternary expression, got %v", expectedElseTok, p.current)
			return ternary
		}

		ternary.ElsePos = p.current.Pos
		ternary.Else = p.current.Kind

		p.next() // Consume elseToken
		p.consumeWhitespaces()

		ternary.FalseExpr = p.parseExpr()

		return ternary
	}

	return condition
}

// parseBinaryExpr parses expressions with operator precedence.
func (p *Parser) parseBinaryExpr(minPrecedence int) Node {
	left := p.parsePrimary()

	for {
		opPrecedence, isRightAssoc := getPrecedence(p.current)
		if opPrecedence < minPrecedence {
			break
		}

		op := p.current
		p.next()

		p.consumeWhitespaces()

		nextMinPrecedence := opPrecedence
		if !isRightAssoc {
			nextMinPrecedence++ // Left-associative operators require higher precedence for right operand
		}

		right := p.parseBinaryExpr(nextMinPrecedence)

		left = &BinaryExpr{
			X:     left,
			Op:    op.Kind,
			OpPos: op.Pos,
			Y:     right,
		}
	}

	return left
}

// parsePrimary handles literals, identifiers, and parenthesized expressions.
func (p *Parser) parsePrimary() Node {
	switch p.current.Kind {
	case token.IDENT:
		ident := Ident{
			Pos:  p.current.Pos,
			Name: p.current.Val,
		}
		p.next()
		p.consumeWhitespaces()
		return ident
	case token.INT, token.FLOAT:
		lit := Lit{
			Pos: p.current.Pos,
			Val: p.current.Val,
			Typ: p.current.Kind,
		}

		p.next()
		p.consumeWhitespaces()
		return lit
	case token.LPAREN:
		p.next() // Consume '('
		expr := p.parseExpr()
		if p.current.Kind != token.RPAREN {
			p.errorf("expected closing ')', got %v", p.current)
			return expr // Return partial expression
		}

		p.next() // Consume ')'
		p.consumeWhitespaces()
		return expr
	default:
		p.errorf("expected identifier, literal, or '(', got %v", p.current)
		return nil
	}
}

// getPrecedence returns the precedence and associativity of an operator.
func getPrecedence(tok token.Token) (int, bool) {
	if tok.IsComparisonOp() {
		return 10, false
	}

	switch tok.Kind {
	case token.ADD, token.SUB:
		return 20, false // Left-associative
	case token.MUL, token.DIV:
		return 30, false // Left-associative
	// case token.POW:
	// 	return 40, true
	default:
		return 0, false
	}
}

func (p *Parser) parseStmt(preBlockWs string) Node {
	if p.current.Kind != token.LSTMT {
		p.errorf("expected LSTMT, got %v", p.current)
		return nil
	}

	p.next() // Consume LSTMT
	p.consumeWhitespaces()

	switch p.current.Kind {
	case token.IF:
		return p.parseIfStmt(preBlockWs)
	case token.GENIF:
		return p.parseGenIfStmt()
	default:
		p.errorf("unexpected statement token: %v", p.current)
		return nil
	}
}

func (p *Parser) parseIfStmt(preBlockWs string) Node {
	ifStmt := IfStmt{
		BegTag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreLStmtWs: preBlockWs,
				LStmt:      p.tokens[p.pos-1].Pos, // Get position of LSTMT
				KwPos:      p.current.Pos,
				Kw:         token.IF,
			},
		},
	}
	p.next() // Consume IF
	p.consumeWhitespaces()

	// Parse the condition
	ifStmt.BegTag.Body = p.parseExpr() // Assuming condition

	if p.current.Kind != token.RSTMT {
		p.errorf("expected RSTMT after if condition, got %v", p.current)
		return ifStmt
	}
	p.next() // Consume RSTMT

	p.consumeLineBreak()

	ifStmt.Body = p.parseBody()

	// Check for "end" statement
	if p.current.Kind != token.LSTMT {
		p.errorf("expected LSTMT for end, got %v", p.current)
		return ifStmt
	}

	ifStmt.EndTag.LStmt = p.current.Pos
	p.next() // Consume LSTMT

	p.consumeWhitespaces()

	if p.current.Kind != token.END {
		p.errorf("expected END, got %v", p.current)
		return ifStmt
	}
	p.next() // Consume END

	p.consumeWhitespaces()

	if p.current.Kind != token.RSTMT {
		p.errorf("expected RSTMT after end, got %v", p.current)
		return ifStmt
	}
	ifStmt.EndTag.RStmt = p.current.Pos
	p.next() // Consume RSTMT

	p.consumeWhitespaces()
	p.consumeLineBreak()

	return ifStmt
}

func (p *Parser) parseGenIfStmt() Node {
	genifStmt := StmtTagWithExpr{
		StmtTag: StmtTag{
			LStmt: p.tokens[p.pos-1].Pos, // Get position of LSTMT
			KwPos: p.current.Pos,
			Kw:    token.GENIF,
		},
	}
	p.next() // Consume GENIF
	p.consumeWhitespaces()

	// Parse the condition
	genifStmt.Body = p.parseExpr() // Assuming condition

	if p.current.Kind != token.RSTMT {
		p.errorf("expected RSTMT after genif condition, got %v", p.current)
		return genifStmt
	}
	p.next() // Consume RSTMT

	p.consumeWhitespaces()
	p.consumeLineBreak()

	return genifStmt
}

func (p *Parser) parseBody() []Node {
	var body []Node
	for {
		if (p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.END) ||
			(p.pos+2 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.WS && p.tokens[p.pos+2].Kind == token.END) {
			break
		}
		node := p.parseNode()
		if node == nil {
			break
		}
		body = append(body, node)
	}
	return body
}
