package token

import (
	"fmt"
	"slices"
)

type Kind int

func (k Kind) String() string {
	return tokens[k]
}

func (k Kind) Bytes() []byte {
	return []byte(tokens[k])
}

func (k Kind) Rune() rune {
	return rune(tokens[k][0])
}

func (k Kind) IsOneOfMany(types ...Kind) bool {
	return slices.Contains(types, k)
}

func (k Kind) IsValueable() bool {
	return valuable_beg < k && k < valuable_end
}

func (k Kind) IsComparasionOp() bool {
	return (comparison_op_beg < k && k < comparison_op_end) ||
		(AND < k && k < ISNOT)
}

func (k Kind) IsLogicalOp() bool {
	return k.IsOneOfMany(AND, LAND, OR, LOR)
}

const (
	EOF Kind = iota

	valuable_beg
	COMM
	LNBR
	TEXT
	WS

	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	STR   // "abc"

	errors_beg // nolint: unused
	// Errors
	NOT_TERMINATED_STR
	errors_end
	valuable_end

	operator_beg
	// Operators and delimiters
	RARR // ->

	// ASSIGN     // =
	// ADD_ASSIGN // +=
	// SUB_ASSIGN // -=
	// MUL_ASSIGN // *=
	// QUO_ASSIGN // /=
	// REM_ASSIGN // %=

	MINUS // -
	ADD   // +
	MUL   // *
	DIV   // /
	MOD   // %

	LEXPR // {{
	REXPR // }}
	LCOMM // {#
	RCOMM // #}
	LSTMT // {%
	RSTMT // %}

	LPAREN // (
	LBRACK // [
	LBRACE // {

	COMMA  // ,
	PERIOD // .

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	comparison_op_beg
	EQL  // ==
	NEQL // !=
	LEQ  // <=
	GEQ  // >=
	LESS // <
	GRTR // >
	LAND // &&
	LOR  // ||
	comparison_op_end

	QUESTION // ?
	COLON    // :
	EXCL     // !

	keyword_beg
	// Keywords

	AND   // and
	OR    // or
	IS    // is
	NOT   // not
	ISNOT // is not

	operator_end
	FOR     // for
	LET     // let
	IF      // if
	ELSE    // else
	SWITCH  // switch
	END     // end
	CASE    // case
	DO      // do
	DEFAULT // default
	EXTEND  // extend
	keyword_end
)

var tokens = []string{
	EOF: "EOF",

	COMM: "COMMENT",
	TEXT: "TEXT",
	LNBR: "LBR",
	WS:   "WHITESPACE",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	STR:   "STRING",

	NOT_TERMINATED_STR: "NOT_TERMINATED_STR",

	QUESTION: "?",
	COLON:    ":",

	MINUS: "-",
	ADD:   "+",
	MUL:   "*",
	DIV:   "/",
	MOD:   "%",

	// ASSIGN:     "=",
	// ADD_ASSIGN: "+=",
	// SUB_ASSIGN: "-=",
	// MUL_ASSIGN: "*=",
	// QUO_ASSIGN: "/=",
	// REM_ASSIGN: "%=",

	EQL:  "==",
	LESS: "<",
	GRTR: ">",
	EXCL: "!",
	NEQL: "!=",
	LEQ:  "<=",
	GEQ:  ">=",
	LAND: "&&",
	LOR:  "||",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",

	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",

	COMMA:  ",",
	PERIOD: ".",

	LEXPR: "{{",
	REXPR: "}}",
	LCOMM: "{#",
	RCOMM: "#}",
	LSTMT: "{%",
	RSTMT: "%}",

	RARR: "->",

	FOR:     "for",
	LET:     "let",
	IF:      "if",
	ELSE:    "else",
	SWITCH:  "switch",
	END:     "end",
	CASE:    "case",
	DEFAULT: "default",
	EXTEND:  "extend",
	AND:     "and",
	OR:      "or",
	IS:      "is",
	NOT:     "not",
	DO:      "do",
	ISNOT:   "is not",
}

const (
	SQUOTE rune = '\''
	DQUOTE rune = '"'
)

type Token struct {
	Kind
	Val string
	Pos Pos
}

func (t Token) String() string {
	if t.IsValueable() {
		switch t.Kind {
		case EOF:
			return "EOF"

		case TEXT:
			return fmt.Sprintf("{Kind: %s, Val: %.10q.., Pos: %s}", t.Kind.String(), t.Val, t.Pos)

		default:
			return fmt.Sprintf("{Kind: %s, Val: %q, Pos: %s}", t.Kind.String(), t.Val, t.Pos)
		}
	}

	return fmt.Sprintf("{Kind: %s, Pos: %s}", t.Kind.String(), t.Pos)
}

func GetOperators() []Kind {
	res := make([]Kind, operator_end-operator_beg)

	for i := range int(operator_end) - int(operator_beg) {
		t := Kind(i + int(operator_beg))

		if t.String() != "" {
			res[i] = t
		}
	}

	return res
}

func GetOperatorsWithoutKw() []Kind {
	res := make([]Kind, keyword_beg-operator_beg)

	for i := range int(keyword_beg) - int(operator_beg) {
		t := Kind(i + int(operator_beg))

		if t.String() != "" {
			res[i] = t
		}
	}

	return res
}

func GetKeywords() []Kind {
	res := make([]Kind, keyword_end-keyword_beg)

	for i := range int(keyword_end) - int(keyword_beg) {
		t := Kind(i + int(keyword_beg))

		if t.String() != "" {
			res[i] = t
		}
	}

	return res
}

func IsNotOp(r rune) bool {
	for i := operator_beg + 1; i < keyword_beg; i++ {
		t := i.String()
		if t != "" && r == rune(t[0]) {
			return false
		}
	}

	return true
}
