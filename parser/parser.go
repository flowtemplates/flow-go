package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

type Parser struct {
	tokens  []token.Token
	pos     int
	nodes   []Node
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

func (p *Parser) Parse() ([]Node, error) {
	for p.pos < len(p.tokens) {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if node != nil {
			p.nodes = append(p.nodes, node)
		} else {
			p.next()
		}
	}

	return p.nodes, nil
}

func (p *Parser) next() token.Token {
	p.pos++
	p.current = p.getCurrent()
	return p.current
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

func (p *Parser) parseNode() (Node, error) {
	switch p.current.Kind {
	case token.TEXT:
		return p.parseText(), nil
	case token.LNBR:
		return p.parseText(), nil
	case token.WS:
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.LSTMT {
			if (p.pos+2 < len(p.tokens) && p.tokens[p.pos+2].Kind == token.END) ||
				(p.pos+3 < len(p.tokens) && p.tokens[p.pos+2].Kind == token.WS && p.tokens[p.pos+3].Kind == token.END) {
				return nil, nil
			}

			ws := p.current.Val
			p.next()
			return p.parseStmt(ws)
		}

		return p.parseText(), nil
	case token.LEXPR:
		return p.parseExprBlock()
	case token.LSTMT:
		return p.parseStmt("")
	// case token.LCOMM:
	// 	return p.parseco()
	// case token.EOF:
	// 	return nil // End of input
	default:
		return nil, errors.New("unexpected token while parsing node")
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
			// TODO: refactor
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.LSTMT {
				if (p.pos+2 < len(p.tokens) && p.tokens[p.pos+2].Kind == token.END) ||
					(p.pos+3 < len(p.tokens) && p.tokens[p.pos+2].Kind == token.WS && p.tokens[p.pos+3].Kind == token.END) {
					break loop
				}
				break loop
			}

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

func (p *Parser) parseExprBlock() (Node, error) {
	exprBlock := ExprBlock{
		LBrace: p.current.Pos,
	}
	p.next() // Consume LEXPR
	p.consumeWhitespaces()

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprBlock.Body = body

	if p.current.Kind != token.REXPR {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.REXPR},
		}
	}

	exprBlock.RBrace = p.current.Pos
	p.next() // Consume REXPR
	return exprBlock, nil
}

func (p *Parser) parseExpr() (Node, error) {
	return p.parseTernaryExpr(1)
}

// parseTernaryExpr parses expressions with ternary operators and precedence.
func (p *Parser) parseTernaryExpr(minPrecedence int) (Node, error) {
	condition, err := p.parseBinaryExpr(minPrecedence)
	if err != nil {
		return nil, err
	}

	doTok := p.current
	if doTok.IsOneOfMany(token.QUESTION, token.DO) {
		ternary := &TernaryExpr{
			Condition: condition,
			DoPos:     p.current.Pos,
			Do:        doTok.Kind,
		}
		p.next() // Consume '?'
		p.consumeWhitespaces()

		trueExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ternary.TrueExpr = trueExpr

		var expectedElseTok token.Kind
		if doTok.Kind == token.QUESTION {
			expectedElseTok = token.COLON
		} else {
			expectedElseTok = token.ELSE
		}

		if p.current.Kind != expectedElseTok {
			return nil, ExpectedTokenError{
				Pos:    p.current.Pos,
				Tokens: []token.Kind{expectedElseTok},
			}
		}

		ternary.ElsePos = p.current.Pos
		ternary.Else = p.current.Kind

		p.next() // Consume elseToken
		p.consumeWhitespaces()

		falseExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ternary.FalseExpr = falseExpr

		return ternary, nil
	}

	return condition, nil
}

func (p *Parser) parseUnaryExpr() (Node, error) {
	if p.current.IsOneOfMany(token.NOT, token.EXCL) {
		op := p.current
		p.next() // Consume operator
		p.consumeWhitespaces()

		operand, err := p.parseUnaryExpr()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Op:    op.Kind,
			OpPos: op.Pos,
			X:     operand,
		}, nil
	}

	// If no unary operator, fallback to primary
	return p.parsePrimary()
}

func (p *Parser) parseBinaryExpr(minPrecedence int) (Node, error) {
	left, err := p.parseUnaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		opPrecedence, isRightAssoc := getPrecedence(p.current)
		if opPrecedence < minPrecedence {
			break
		}

		op := p.current
		p.next()

		p.consumeWhitespaces()
		if op.Kind == token.IS && p.current.Kind == token.NOT {
			op = token.Token{
				Kind: token.ISNOT,
				Pos:  op.Pos,
			}
			p.next()
			p.consumeWhitespaces()
		}

		nextMinPrecedence := opPrecedence
		if !isRightAssoc {
			nextMinPrecedence++ // Left-associative operators require higher precedence for right operand
		}

		right, err := p.parseBinaryExpr(nextMinPrecedence)
		if err != nil {
			return nil, err
		}

		left = &BinaryExpr{
			X:     left,
			Op:    op.Kind,
			OpPos: op.Pos,
			Y:     right,
		}
	}

	return left, nil
}

// parsePrimary handles literals, identifiers, and parenthesized expressions.
func (p *Parser) parsePrimary() (Node, error) {
	switch p.current.Kind {
	case token.IDENT:
		ident := Ident{
			Pos:  p.current.Pos,
			Name: p.current.Val,
		}
		p.next()
		p.consumeWhitespaces()
		return ident, nil
	case token.INT, token.FLOAT, token.STR:
		var val value.Valueable
		switch p.current.Kind {
		case token.INT, token.FLOAT:
			v, err := strconv.ParseFloat(p.current.Val, 64)
			if err != nil {
				panic(err)
			}

			val = value.NumberValue(v)
		case token.STR:
			// TODO: check for quotes
			val = value.StringValue(p.current.Val[1 : len(p.current.Val)-1])
		}

		lit := Lit{
			Pos:   p.current.Pos,
			Value: val,
		}

		p.next()
		p.consumeWhitespaces()
		return lit, nil
	case token.LPAREN:
		p.next() // Consume '('
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.current.Kind != token.RPAREN {
			return nil, ExpectedTokenError{
				Pos:    p.current.Pos,
				Tokens: []token.Kind{token.RPAREN},
			}
		}

		p.next() // Consume ')'
		p.consumeWhitespaces()
		return expr, nil
	default:
		return nil, Error{
			Pos: p.current.Pos,
			Typ: ErrExpressionExpected,
		}
	}
}

func getPrecedence(tok token.Token) (int, bool) {
	if tok.IsComparisonOp() {
		return 10, false
	}

	switch tok.Kind {
	case token.LOR, token.OR:
		return 10, false
	case token.AND, token.LAND:
		return 20, false
	// case token.ADD, token.SUB:
	// 	return 20, false // Left-associative
	// case token.MUL, token.DIV:
	// 	return 30, false // Left-associative
	// case token.POW:
	// 	return 40, true
	default:
		return 0, false
	}
}

func (p *Parser) parseStmt(preBlockWs string) (Node, error) {
	if p.current.Kind != token.LSTMT {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.LSTMT},
		}
	}

	p.next() // Consume LSTMT
	p.consumeWhitespaces()

	switch p.current.Kind {
	case token.IF:
		return p.parseIfStmt(preBlockWs)
	case token.GENIF:
		return p.parseGenIfStmt()
	default:
		return nil, errors.New("")
	}
}

func (p *Parser) parseIfStmt(preBlockWs string) (Node, error) {
	ifStmt := IfStmt{
		BegTag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: preBlockWs,
				LStmt: p.tokens[p.pos-1].Pos, // Get position of LSTMT
				KwPos: p.current.Pos,
				Kw:    token.IF,
			},
		},
	}
	p.next() // Consume IF
	p.consumeWhitespaces()

	begTagBody, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ifStmt.BegTag.Body = begTagBody

	if p.current.Kind != token.RSTMT {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeLineBreak()

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	ifStmt.Body = body

	ifStmt.PreEndTagWs = p.consumeWhitespaces()
	// Check for "end" statement
	if p.current.Kind != token.LSTMT {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.LSTMT},
		}
	}

	p.next() // Consume LSTMT

	p.consumeWhitespaces()

	if p.current.Kind != token.END {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.END},
		}
	}
	p.next() // Consume END

	p.consumeWhitespaces()

	if p.current.Kind != token.RSTMT {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespaces()
	p.consumeLineBreak()

	return ifStmt, nil
}

func (p *Parser) parseGenIfStmt() (Node, error) {
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

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	genifStmt.Body = body

	if p.current.Kind != token.RSTMT {
		return nil, ExpectedTokenError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespaces()
	p.consumeLineBreak()

	return genifStmt, nil
}

func (p *Parser) parseBody() ([]Node, error) {
	var body []Node
	for {
		if (p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.END) ||
			(p.pos+2 < len(p.tokens) && p.tokens[p.pos+1].Kind == token.WS && p.tokens[p.pos+2].Kind == token.END) {
			break
		}

		node, err := p.parseNode()
		if err != nil {
			return body, err
		}

		if node == nil {
			fmt.Println("12312312312312")
			break
		}
		body = append(body, node)
	}
	return body, nil
}
