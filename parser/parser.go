package parser

import (
	"errors"
	"strings"
	"unicode"

	"github.com/flowtemplates/flow-go/token"
)

type parser struct {
	tokens       []token.Token
	pos          int
	currentToken token.Token
}

func newParser(tokens []token.Token) *parser {
	p := &parser{
		tokens: tokens,
		pos:    -1,
	}

	p.next()

	return p
}

func (p *parser) parse() (Ast, error) {
	var nodes []Node

	for p.pos < len(p.tokens) {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if node != nil {
			nodes = append(nodes, node)
		} else {
			p.next()
		}
	}

	return nodes, nil
}

func (p *parser) getCurrent() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}

	return token.Token{Kind: token.EOF}
}

func (p *parser) next() {
	p.pos++
	p.currentToken = p.getCurrent()
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

func trimSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) && r != '\n' && r != '\r'
	})
}

func (p *parser) parseComm(preWs string) (*CommNode, error) {
	commNode := CommNode{
		Pos:   p.currentToken.Pos,
		PreWs: preWs,
	}

	p.next() // Consume LCOMM

	switch p.currentToken.Kind {
	case token.COMM_TEXT:
		commNode.Val = trimSpaces(p.currentToken.Val)

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

	case token.SWITCH:
		return p.parseSwitchStmt(preWs)

	default:
		return nil, Error{
			Pos: p.currentToken.Pos,
			Typ: ErrKeywordExpected,
		}
	}
}

func (p *parser) consumeEndTag() error {
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

func (p *parser) parseElses(ifStmt *IfNode) error {
	for {
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

			ifStmt.EndTag = StmtTag{PreWs: preTagWs}

			return nil

		case token.ELSE:
			p.next()
			p.consumeWhitespace()

			if p.currentToken.Kind == token.IF {
				// Found "else if"
				elseIfNode, err := p.parseElseIf(preTagWs)
				if err != nil {
					return err
				}

				ifStmt.ElseIfs = append(ifStmt.ElseIfs, elseIfNode)

				continue
			}

			// Standard "else"
			if p.currentToken.Kind != token.RSTMT {
				return ExpectedTokensError{
					Pos:    p.currentToken.Pos,
					Tokens: []token.Kind{token.RSTMT},
				}
			}

			p.next() // Consume RSTMT

			p.consumeLineBreak()

			elseBody, err := p.parseBody()
			if err != nil {
				return err
			}

			ifStmt.Else = Clause{
				Tag:  StmtTag{PreWs: preTagWs},
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

			ifStmt.EndTag = StmtTag{PreWs: preEndTagWs}

			return nil

		default:
			return Error{
				Pos: p.currentToken.Pos,
				Typ: ErrEndExpected,
			}
		}
	}
}

func (p *parser) parseElseIf(preTagWs string) (ClauseWithExpr, error) {
	p.next() // Consume IF
	p.consumeWhitespace()

	expr, err := p.parseExpr()
	if err != nil {
		return ClauseWithExpr{}, err
	}

	if p.currentToken.Kind != token.RSTMT {
		return ClauseWithExpr{}, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}

	p.next() // Consume RSTMT

	p.consumeLineBreak()

	body, err := p.parseBody()
	if err != nil {
		return ClauseWithExpr{}, err
	}

	return ClauseWithExpr{
		Tag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: preTagWs,
			},
			Expr: expr,
		},
		Body: body,
	}, nil
}

func (p *parser) parseIfStmt(preWs string) (Node, error) {
	ifStmt := IfNode{
		IfTag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: preWs,
			},
		},
	}

	p.next() // Consume IF
	p.consumeWhitespace()

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

	p.consumeWhitespace()
	p.consumeLineBreak()

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}

	ifStmt.Main = body

	if err := p.parseElses(&ifStmt); err != nil {
		return nil, err
	}

	return &ifStmt, nil
}

func (p *parser) parseGenIfStmt() (Node, error) {
	genifStmt := GenifNode{
		StmtTagWithExpr: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: "", // TODO: fix
			},
		},
	}

	p.next() // Consume GENIF
	p.consumeWhitespace()

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

	for p.currentToken.Kind != token.LSTMT ||
		(!p.checkNextNTokens(token.END) &&
			!p.checkNextNTokens(token.WS, token.END) &&
			!p.checkNextNTokens(token.ELSE) &&
			!p.checkNextNTokens(token.WS, token.ELSE) &&
			!p.checkNextNTokens(token.CASE) &&
			!p.checkNextNTokens(token.WS, token.CASE) &&
			!p.checkNextNTokens(token.DEFAULT) &&
			!p.checkNextNTokens(token.WS, token.DEFAULT)) {
		// TODO: refactor
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

func (p *parser) parseCases(switchStmt *SwitchNode) error {
	for {
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

			switchStmt.EndTag = StmtTag{PreWs: preTagWs}

			return nil

		case token.CASE:
			cc := ClauseWithExpr{
				Tag: StmtTagWithExpr{
					StmtTag: StmtTag{
						PreWs: preTagWs,
					},
				},
			}

			p.next()
			p.consumeWhitespace()

			cExpr, err := p.parseExpr()
			if err != nil {
				return err
			}

			cc.Tag.Expr = cExpr

			if p.currentToken.Kind != token.RSTMT {
				return ExpectedTokensError{
					Pos:    p.currentToken.Pos,
					Tokens: []token.Kind{token.RSTMT},
				}
			}

			p.next() // Consume RSTMT

			p.consumeLineBreak()

			b, err := p.parseBody()
			if err != nil {
				return err
			}

			cc.Body = b

			switchStmt.Cases = append(switchStmt.Cases, cc)

		case token.DEFAULT:
			d := Clause{
				Tag: StmtTag{
					PreWs: preTagWs,
				},
			}

			p.next()
			p.consumeWhitespace()

			if p.currentToken.Kind != token.RSTMT {
				return ExpectedTokensError{
					Pos:    p.currentToken.Pos,
					Tokens: []token.Kind{token.RSTMT},
				}
			}

			p.next() // Consume RSTMT

			p.consumeLineBreak()

			b, err := p.parseBody()
			if err != nil {
				return err
			}

			d.Body = b

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

			switchStmt.EndTag = StmtTag{PreWs: preEndTagWs}
			switchStmt.DefaultCase = &d

			return nil

		default:
			return Error{
				Pos: p.currentToken.Pos,
				Typ: ErrEndExpected,
			}
		}
	}
}

func (p *parser) parseSwitchStmt(preWs string) (Node, error) {
	switchStmt := SwitchNode{
		SwitchTag: StmtTagWithExpr{
			StmtTag: StmtTag{
				PreWs: preWs,
			},
		},
	}

	p.next() // Consume SWITCH
	p.consumeWhitespace()

	begTagBody, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	switchStmt.SwitchTag.Expr = begTagBody

	if p.currentToken.Kind != token.RSTMT {
		return nil, ExpectedTokensError{
			Pos:    p.currentToken.Pos,
			Tokens: []token.Kind{token.RSTMT},
		}
	}

	p.next() // Consume RSTMT

	p.consumeWhitespace()
	p.consumeLineBreak()

	if err := p.parseCases(&switchStmt); err != nil {
		return nil, err
	}

	return &switchStmt, nil
}
