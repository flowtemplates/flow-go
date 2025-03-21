package parser

import (
	"github.com/flowtemplates/flow-go/token"
)

func (p *parser) next() token.Token {
	p.pos++
	p.current = p.getCurrent()
	return p.current
}

func (p *parser) consumeToken(t token.Kind) string {
	var val string
	if p.current.Kind == t {
		val = p.current.Val
		p.next()
	}
	return val
}

func (p *parser) consumeLineBreak() string {
	return p.consumeToken(token.LNBR)
}

func (p *parser) consumeWhitespace() string {
	return p.consumeToken(token.WS)
}

func (p *parser) checkNextNTokens(tokens ...token.Kind) bool {
	if p.pos+len(tokens) > len(p.tokens) {
		return false
	}

	for i, kind := range tokens {
		if p.tokens[p.pos+i+1].Kind != kind {
			return false
		}
	}

	return true
}

func (p *parser) getCurrent() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Kind: token.EOF}
}

func (p *parser) parseNode() (Node, error) {
	switch p.current.Kind {
	case token.TEXT:
		return p.parseText(), nil
	case token.LNBR:
		return p.parseText(), nil
	case token.WS:
		if p.checkNextNTokens(token.LSTMT) {
			if p.checkNextNTokens(token.LSTMT, token.END) || p.checkNextNTokens(token.LSTMT, token.WS, token.END) {
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
		return nil, nil
	}
}

func (p *parser) parseText() Node {
	var res []string
	for p.current.IsOneOfMany(token.TEXT, token.LNBR, token.WS) {
		if p.current.Kind == token.WS && p.checkNextNTokens(token.LSTMT) {
			break
		}

		res = append(res, p.current.Val)
		p.next()
	}

	text := Text{
		Pos: p.current.Pos,
		Val: res,
	}
	return text
}

func (p *parser) parseStmt(preBlockWs string) (Node, error) {
	if p.current.Kind != token.LSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.LSTMT},
		}
	}

	p.next() // Consume LSTMT
	p.consumeWhitespace()

	switch p.current.Kind {
	case token.IF:
		return p.parseIfStmt(preBlockWs)
	case token.GENIF:
		return p.parseGenIfStmt()
	// case token.SWITCH:
	// 	return p.parseSwitchStmt(preBlockWs)
	default:
		return nil, Error{
			Pos: p.current.Pos,
			Typ: ErrKeywordExpected,
		}
	}
}

func (p *parser) consumeEndTag() (string, error) {
	preEndTagWs := p.consumeWhitespace()
	// Check for "end" statement
	if p.current.Kind != token.LSTMT {
		return "", Error{
			Pos: p.current.Pos,
			Typ: ErrEndExpected,
		}
	}

	p.next() // Consume LSTMT

	p.consumeWhitespace()

	if p.current.Kind != token.END {
		return "", ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.END},
		}
	}
	p.next() // Consume END

	p.consumeWhitespace()

	if p.current.Kind != token.RSTMT {
		return "", ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespace()
	p.consumeLineBreak()

	return preEndTagWs, nil
}

func (p *parser) parseIfStmt(preBlockWs string) (Node, error) {
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
	p.consumeWhitespace()

	begTagBody, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ifStmt.BegTag.Body = begTagBody

	if p.current.Kind != token.RSTMT {
		return nil, ExpectedTokensError{
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

	preEndTagWs, err := p.consumeEndTag()
	if err != nil {
		return nil, err
	}
	ifStmt.PreEndTagWs = preEndTagWs

	return ifStmt, nil
}

// func (p *Parser) parseSwitchStmt(preBlockWs string) (Node, error) {
// 	ifStmt := SwitchStmt{
// 		BegTag: StmtTagWithExpr{
// 			StmtTag: StmtTag{
// 				PreWs: preBlockWs,
// 				LStmt: p.tokens[p.pos-1].Pos,
// 				KwPos: p.current.Pos,
// 				Kw:    token.SWITCH,
// 			},
// 		},
// 	}
// 	p.next() // Consume SWITCH
// 	p.consumeWhitespaces()

// 	begTagBody, err := p.parseExpr()
// 	if err != nil {
// 		return nil, err
// 	}
// 	ifStmt.BegTag.Body = begTagBody

// 	if p.current.Kind != token.RSTMT {
// 		return nil, ExpectedTokenError{
// 			Pos:    p.current.Pos,
// 			Tokens: []token.Kind{token.RSTMT},
// 		}
// 	}
// 	p.next() // Consume RSTMT

// 	p.consumeLineBreak()

// 	body, err := p.parseBody()
// 	if err != nil {
// 		return nil, err
// 	}
// 	ifStmt.Body = body

// 	ifStmt.PreEndTagWs = p.consumeWhitespaces()
// 	// Check for "end" statement
// 	if p.current.Kind != token.LSTMT {
// 		return nil, ExpectedTokenError{
// 			Pos:    p.current.Pos,
// 			Tokens: []token.Kind{token.LSTMT},
// 		}
// 	}

// 	p.next() // Consume LSTMT

// 	p.consumeWhitespaces()

// 	if p.current.Kind != token.END {
// 		return nil, ExpectedTokenError{
// 			Pos:    p.current.Pos,
// 			Tokens: []token.Kind{token.END},
// 		}
// 	}
// 	p.next() // Consume END

// 	p.consumeWhitespaces()

// 	if p.current.Kind != token.RSTMT {
// 		return nil, ExpectedTokenError{
// 			Pos:    p.current.Pos,
// 			Tokens: []token.Kind{token.RSTMT},
// 		}
// 	}
// 	p.next() // Consume RSTMT

// 	p.consumeWhitespaces()
// 	p.consumeLineBreak()

// 	return ifStmt, nil
// }

func (p *parser) parseGenIfStmt() (Node, error) {
	genifStmt := StmtTagWithExpr{
		StmtTag: StmtTag{
			LStmt: p.tokens[p.pos-1].Pos, // Get position of LSTMT
			KwPos: p.current.Pos,
			Kw:    token.GENIF,
		},
	}
	p.next() // Consume GENIF
	p.consumeWhitespace()

	// Parse the condition

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	genifStmt.Body = body

	if p.current.Kind != token.RSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.current.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespace()
	p.consumeLineBreak()

	return genifStmt, nil
}

func (p *parser) parseBody() ([]Node, error) {
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
			break
		}
		body = append(body, node)
	}
	return body, nil
}
