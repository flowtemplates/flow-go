package parser

import (
	"errors"

	"github.com/flowtemplates/flow-go/token"
)

func (p *parser) getCurrent() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Kind: token.EOF}
}

func (p *parser) next() token.Token {
	p.pos++
	p.currentToken = p.getCurrent()
	return p.currentToken
}

func (p *parser) consumeToken(t token.Kind) string {
	var val string
	if p.currentToken.Kind == t {
		val = p.currentToken.Val
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
	if p.pos+len(tokens)+1 > len(p.tokens) {
		return false
	}

	for i, kind := range tokens {
		if p.tokens[p.pos+i+1].Kind != kind {
			return false
		}
	}

	return true
}

func (p *parser) parseNode() (Node, error) {
	switch p.currentToken.Kind {
	case token.TEXT:
		return p.parseText(), nil
	case token.LNBR:
		return p.parseText(), nil
	case token.WS:
		if p.checkNextNTokens(token.LSTMT) {
			if p.checkNextNTokens(token.LSTMT, token.END) || p.checkNextNTokens(token.LSTMT, token.WS, token.END) ||
				p.checkNextNTokens(token.LSTMT, token.ELSE) || p.checkNextNTokens(token.LSTMT, token.WS, token.ELSE) {
				return nil, nil
			}

			ws := p.currentToken.Val
			p.next()
			return p.parseStmt(ws)
		}

		if p.checkNextNTokens(token.LCOMM) {
			// if p.checkNextNTokens(token.LSTMT, token.END) || p.checkNextNTokens(token.LSTMT, token.WS, token.END) ||
			// 	p.checkNextNTokens(token.LSTMT, token.ELSE) || p.checkNextNTokens(token.LSTMT, token.WS, token.ELSE) {
			// 	return nil, nil
			// }

			ws := p.currentToken.Val
			p.next()
			return p.parseComm(ws)
		}

		return p.parseText(), nil
	case token.LEXPR:
		return p.parseExprNode()
	case token.LSTMT:
		return p.parseStmt("")
	case token.LCOMM:
		return p.parseComm("")
	// case token.EOF:
	// 	return nil // End of input
	default:
		return nil, nil
	}
}

func (p *parser) parseText() *TextNode {
	var res []string
	for p.currentToken.IsOneOfMany(token.TEXT, token.LNBR, token.WS) {
		if p.currentToken.Kind == token.WS && p.checkNextNTokens(token.LSTMT) {
			break
		}

		res = append(res, p.currentToken.Val)
		p.next()
	}

	text := TextNode{
		Pos: p.currentToken.Pos,
		Val: res,
	}

	return &text
}

func (p *parser) parseComm(preWs string) (*CommNode, error) {
	commNode := CommNode{
		Pos:   p.currentToken.Pos,
		PreWs: preWs,
	}
	p.next() // Consume LCOMM
	switch p.currentToken.Kind {
	case token.COMM_TEXT:
		commNode.Val = p.currentToken.Val
		p.next()
	case token.RCOMM:

	default:
		return nil, errors.New("unexpected token inside comment")
	}
	if p.currentToken.Kind != token.RCOMM {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RCOMM},
		}
	}

	p.next()
	p.consumeWhitespace()
	commNode.PostLB = p.consumeLineBreak()

	return &commNode, nil
}

func (p *parser) parseStmt(preWs string) (Node, error) {
	if p.currentToken.Kind != token.LSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.LSTMT},
		}
	}

	p.next() // Consume LSTMT
	p.consumeWhitespace()

	switch p.currentToken.Kind {
	case token.IF:
		return p.parseIfStmt(preWs)
	case token.GENIF:
		return p.parseGenIfStmt()
	// case token.SWITCH:
	// 	return p.parseSwitchStmt(preBlockWs)
	default:
		return nil, Error{
			Pos: p.currentToken.Pos,
			Typ: ErrKeywordExpected,
		}
	}
}

func (p *parser) consumeEndTag() error {
	// if p.currentToken.Kind != token.LSTMT {
	// 	return Error{
	// 		Pos: p.currentToken.Pos,
	// 		Typ: ErrEndExpected,
	// 	}
	// }

	// p.next() // Consume LSTMT

	// p.consumeWhitespace()

	// if p.currentToken.Kind != token.END {
	// 	return ExpectedTokensError{
	// 		Pos:    p.currentToken.Pos,
	// 		Tokens: []token.Kind{token.END},
	// 	}
	// }
	p.next() // Consume END

	p.consumeWhitespace()

	if p.currentToken.Kind != token.RSTMT {
		return ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespace()
	p.consumeLineBreak()

	return nil
}

func (p *parser) parseIfStmt(preWs string) (Node, error) {
	ifStmt := IfNode{
		IfTag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: preWs,
				// LStmt: p.tokens[p.pos-1].Pos, // Get position of LSTMT
				// KwPos: p.currentToken.Pos,
			},
		},
	}
	p.next() // Consume IF
	p.consumeWhitespace()

	// startPos := p.pos
	// for !p.currentToken.IsOneOfMany(token.RSTMT, token.LNBR) {
	// 	p.next()
	// }

	begTagBody, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	ifStmt.IfTag.Expr = begTagBody

	if p.currentToken.Kind != token.RSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeLineBreak()

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	ifStmt.MainBody = body

	if err := p.parseElses(&ifStmt); err != nil {
		return nil, err
	}

	return &ifStmt, nil
}

func (p *parser) parseElses(ifStmt *IfNode) error {
	preTagWs := p.consumeWhitespace()

	if p.currentToken.Kind != token.LSTMT {
		return Error{
			Pos: p.currentToken.Pos,
			Typ: ErrEndExpected,
		}
	}

	p.next() // Consume LSTMT

	p.consumeWhitespace()

	switch p.currentToken.Kind {
	case token.END:
		if err := p.consumeEndTag(); err != nil {
			return err
		}
		ifStmt.EndTag = StmtTag{
			PreWs: preTagWs,
		}
	case token.ELSE:
		p.next()
		p.consumeWhitespace()
		switch p.currentToken.Kind {
		case token.RSTMT:
			p.next() // Consume RSTMT

			p.consumeLineBreak()

			elseBody, err := p.parseBody()
			if err != nil {
				return err
			}

			ifStmt.ElseBody = ElseNode{
				ElseTag: StmtTag{
					PreWs: preTagWs,
				},
				Body: elseBody,
			}

			preEndTagWs := p.consumeWhitespace()

			if p.currentToken.Kind != token.LSTMT {
				return Error{
					Pos: p.currentToken.Pos,
					Typ: ErrEndExpected,
				}
			}

			p.next() // Consume LSTMT

			p.consumeWhitespace()

			if err := p.consumeEndTag(); err != nil {
				return err
			}
			ifStmt.EndTag = StmtTag{
				PreWs: preEndTagWs,
			}
		}
	default:
		return Error{
			Pos: p.currentToken.Pos,
			Typ: ErrEndExpected,
		}
	}

	return nil
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
	genifStmt := StmtNode{
		StmtTagWithKw: StmtTagWithKw{
			Kw: Kw{
				Kind: token.GENIF,
				Pos:  p.currentToken.Pos,
			},
			StmtTag: StmtTag{
				PreWs: "", // TODO: fix
			},
		},
	}
	p.next() // Consume GENIF
	p.consumeWhitespace()

	// Parse the condition

	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	genifStmt.Expr = body

	if p.currentToken.Kind != token.RSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}
	p.next() // Consume RSTMT

	p.consumeWhitespace()
	p.consumeLineBreak()

	return &genifStmt, nil
}

func (p *parser) parseBody() ([]Node, error) {
	var body []Node
	for {
		if p.currentToken.Kind == token.LSTMT &&
			(p.checkNextNTokens(token.END) || p.checkNextNTokens(token.WS, token.END) ||
				p.checkNextNTokens(token.ELSE) || p.checkNextNTokens(token.WS, token.ELSE)) {
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
