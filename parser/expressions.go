package parser

import (
	"strconv"

	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

func (p *parser) parseExprNode() (*ExprNode, error) {
	exprNode := ExprNode{}
	p.next() // Consume LEXPR
	p.consumeWhitespace()

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprNode.Body = body

	if p.currentToken.Kind != token.REXPR {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.REXPR},
		}
	}

	p.next() // Consume REXPR
	return &exprNode, nil
}

func (p *parser) parseExpr() (Expr, error) {
	return p.parseFilterExpr()
}

func (p *parser) parseFilterExpr() (Expr, error) {
	expr, err := p.parseTernaryExpr(1)
	if err != nil {
		return nil, err
	}

	for p.currentToken.Kind == token.RARR {
		opPos := p.currentToken.Pos
		p.next()
		p.consumeWhitespace()

		if p.currentToken.Kind != token.IDENT {
			return nil, ExpectedTokensError{
				Pos:    p.currentToken.Pos,
				Tokens: []token.Kind{token.IDENT},
			}
		}

		ident := Ident{
			Pos:  p.currentToken.Pos,
			Name: p.currentToken.Val,
		}
		p.next()
		p.consumeWhitespace()

		expr = &FilterExpr{
			Expr:   expr,
			OpPos:  opPos,
			Filter: ident,
		}
	}

	return expr, nil
}

func (p *parser) parseTernaryExpr(minPrecedence int) (Expr, error) {
	condition, err := p.parseBinaryExpr(minPrecedence)
	if err != nil {
		return nil, err
	}

	doTok := p.currentToken
	if doTok.IsOneOfMany(token.QUESTION, token.DO) {
		ternary := &TernaryExpr{
			Condition: condition,
			Do: Kw{
				Kind: doTok.Kind,
				Pos:  p.currentToken.Pos,
			},
		}
		p.next() // Consume '?'
		p.consumeWhitespace()

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

		if p.currentToken.Kind != expectedElseTok {
			return nil, ExpectedTokensError{
				Pos:    p.currentToken.Pos,
				Tokens: []token.Kind{expectedElseTok},
			}
		}

		ternary.Else = Kw{
			Kind: p.currentToken.Kind,
			Pos:  p.currentToken.Pos,
		}

		p.next() // Consume elseToken
		p.consumeWhitespace()

		falseExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ternary.FalseExpr = falseExpr

		return ternary, nil
	}

	return condition, nil
}

func (p *parser) parseUnaryExpr() (Expr, error) {
	if p.currentToken.IsOneOfMany(token.NOT, token.EXCL) {
		op := p.currentToken
		p.next() // Consume operator
		p.consumeWhitespace()

		operand, err := p.parseUnaryExpr()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Op: Kw{
				Kind: op.Kind,
				Pos:  op.Pos,
			},
			Expr: operand,
		}, nil
	}

	return p.parsePrimary()
}

func (p *parser) parseBinaryExpr(minPrecedence int) (Expr, error) {
	left, err := p.parseUnaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		opPrecedence, isRightAssoc := getPrecedence(p.currentToken)
		if opPrecedence < minPrecedence {
			break
		}

		op := p.currentToken
		p.next()

		p.consumeWhitespace()
		if op.Kind == token.IS && p.currentToken.Kind == token.NOT {
			op = token.Token{
				Kind: token.ISNOT,
				Pos:  op.Pos,
			}
			p.next()
			p.consumeWhitespace()
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
			X: left,
			Op: Kw{
				Kind: op.Kind,
				Pos:  op.Pos,
			},
			Y: right,
		}
	}

	return left, nil
}

// parsePrimary handles literals, identifiers, and parenthesized expressions.
func (p *parser) parsePrimary() (Expr, error) {
	switch p.currentToken.Kind {
	case token.IDENT:
		ident := Ident{
			Pos:  p.currentToken.Pos,
			Name: p.currentToken.Val,
		}
		p.next()
		p.consumeWhitespace()
		return &ident, nil
	case token.STR:
		lit := &StringLit{
			Pos:   p.currentToken.Pos,
			Quote: p.currentToken.Val[0],
			Value: value.StringValue(p.currentToken.Val[1 : len(p.currentToken.Val)-1]),
		}

		p.next()
		p.consumeWhitespace()
		return lit, nil
	case token.MINUS, token.INT, token.FLOAT:
		var negative bool
		if p.currentToken.Kind == token.MINUS {
			p.next()
			negative = true
		}

		var lit Expr
		switch p.currentToken.Kind {
		case token.INT, token.FLOAT:
			v, err := strconv.ParseFloat(p.currentToken.Val, 64)
			if err != nil {
				panic(err)
			}

			if negative {
				v *= -1
			}

			lit = &NumberLit{
				Pos:   p.currentToken.Pos,
				Value: value.NumberValue(v),
			}
		case token.STR:
			quote := p.currentToken.Val[0]

			lit = &StringLit{
				Pos:   p.currentToken.Pos,
				Quote: quote,
				Value: value.StringValue(p.currentToken.Val[1 : len(p.currentToken.Val)-1]),
			}
		}

		p.next()
		p.consumeWhitespace()
		return lit, nil
	case token.LPAREN:
		parenExpr := ParenExpr{
			Lparen: p.currentToken.Pos,
		}
		p.next() // Consume '('

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.currentToken.Kind != token.RPAREN {
			return nil, ExpectedTokensError{
				Pos:    p.currentToken.Pos,
				Tokens: []token.Kind{token.RPAREN},
			}
		}

		parenExpr.Rparen = p.currentToken.Pos
		p.next()
		p.consumeWhitespace()
		parenExpr.Expr = expr

		return &parenExpr, nil
	// case token.REXPR:
	default:
		return nil, Error{
			Pos: p.currentToken.Pos,
			Typ: ErrExpressionExpected,
		}
		// return nil, ExpectedTokensError{
		// 	Pos:    p.currentToken.Pos,
		// 	Tokens: []token.Kind{token.REXPR},
		// }
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
