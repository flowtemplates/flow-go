package parser

import (
	"strconv"

	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

func (p *parser) parseExprBlock() (Node, error) {
	exprBlock := ExprBlock{
		LBrace: p.current.Pos,
	}
	p.next() // Consume LEXPR
	p.consumeWhitespace()

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprBlock.Body = body

	if p.current.Kind != token.REXPR {
		return nil, ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.REXPR},
		}
	}

	exprBlock.RBrace = p.current.Pos
	p.next() // Consume REXPR
	return exprBlock, nil
}

func (p *parser) parseExpr() (Node, error) {
	return p.parseTernaryExpr(1)
}

// parseTernaryExpr parses expressions with ternary operators and precedence.
func (p *parser) parseTernaryExpr(minPrecedence int) (Node, error) {
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

		if p.current.Kind != expectedElseTok {
			return nil, ExpectedTokensError{
				Pos:    p.current.Pos,
				Tokens: []token.Kind{expectedElseTok},
			}
		}

		ternary.ElsePos = p.current.Pos
		ternary.Else = p.current.Kind

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

func (p *parser) parseUnaryExpr() (Node, error) {
	if p.current.IsOneOfMany(token.NOT, token.EXCL) {
		op := p.current
		p.next() // Consume operator
		p.consumeWhitespace()

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

func (p *parser) parseBinaryExpr(minPrecedence int) (Node, error) {
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

		p.consumeWhitespace()
		if op.Kind == token.IS && p.current.Kind == token.NOT {
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
			X:     left,
			Op:    op.Kind,
			OpPos: op.Pos,
			Y:     right,
		}
	}

	return left, nil
}

// parsePrimary handles literals, identifiers, and parenthesized expressions.
func (p *parser) parsePrimary() (Node, error) {
	switch p.current.Kind {
	case token.IDENT:
		ident := Ident{
			Pos:  p.current.Pos,
			Name: p.current.Val,
		}
		p.next()
		p.consumeWhitespace()
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
		p.consumeWhitespace()
		return lit, nil
	case token.LPAREN:
		p.next() // Consume '('
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.current.Kind != token.RPAREN {
			return nil, ExpectedTokensError{
				Pos:    p.current.Pos,
				Tokens: []token.Kind{token.RPAREN},
			}
		}

		p.next() // Consume ')'
		p.consumeWhitespace()
		return expr, nil
	case token.REXPR:
		return nil, Error{
			Pos: p.current.Pos,
			Typ: ErrExpressionExpected,
		}
	default:
		return nil, ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.REXPR},
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
