package token

import (
	"fmt"
	"slices"
)

type Kind int

func (k Kind) String() string {
	return TokenString(k)
}

const (
	EOF Kind = iota
	ILLEGAL

	valueable_beg
	COMM_TEXT
	LNBR
	TEXT
	WS

	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	STR   // "abc"

	errors_beg
	// Errors
	NOT_TERMINATED_STR
	errors_end
	valueable_end

	operator_beg
	// Operators and delimiters
	RARR // ->

	LAND // &&
	LOR  // ||
	comparison_op_beg
	EQL  // ==
	NEQL // !=
	LEQ  // <=
	GEQ  // >=
	LESS // <
	GRTR // >
	comparison_op_end

	// ASSIGN     // =
	// ADD_ASSIGN // +=
	// SUB_ASSIGN // -=
	// MUL_ASSIGN // *=
	// QUO_ASSIGN // /=
	// REM_ASSIGN // %=

	MINUS // -
	// ADD // +
	// MUL // *
	// DIV // /
	// MOD // %

	QUESTION // ?
	COLON    // :
	EXCL     // !

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
	operator_end

	keyword_beg
	// Keywords
	FOR     // for
	LET     // let
	IF      // if
	GENIF   // genif
	ELSE    // else
	SWITCH  // switch
	END     // end
	CASE    // case
	DEFAULT // default
	EXTEND  // extend
	AND     // and
	OR      // or
	DO      // do
	IS      // is
	NOT     // not
	ISNOT   // is not
	keyword_end
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	COMM_TEXT: "COMMENT",
	TEXT:      "TEXT",
	LNBR:      "LBR",
	WS:        "WHITESPACE",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	STR:   "STRING",

	NOT_TERMINATED_STR: "NOT_TERMINATED_STR",

	QUESTION: "?",
	COLON:    ":",

	MINUS: "-",
	// ADD: "+",
	// MUL: "*",
	// DIV: "/",
	// MOD: "%",

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
	GENIF:   "genif",
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

func TokenString(k Kind) string {
	return tokens[k]
}

func TokenRune(k Kind) rune {
	return rune(tokens[k][0])
}

type Token struct {
	Kind Kind
	Val  string
	Pos  Position
}

func (t Token) String() string {
	if t.IsValueable() {
		switch t.Kind {
		case EOF:
			return "EOF"
		case TEXT:
			return fmt.Sprintf("{Kind: %s, Val: %.10q, Pos: %s}", TokenString(t.Kind), t.Val, t.Pos)
		default:
			return fmt.Sprintf("{Kind: %s, Val: %q, Pos: %s}", TokenString(t.Kind), t.Val, t.Pos)
		}
	}

	return fmt.Sprintf("{Kind: %s, Pos: %s}", TokenString(t.Kind), t.Pos)
}

func (t Token) IsOneOfMany(types ...Kind) bool {
	return slices.Contains(types, t.Kind)
}

func (t Token) IsValueable() bool {
	return valueable_beg < t.Kind && t.Kind < valueable_end
}

func GetOperators() []Kind {
	res := make([]Kind, operator_end-operator_beg)

	for i := range int(operator_end) - int(operator_beg) {
		res[i] = Kind(i + int(operator_beg))
	}

	return res
}

func GetKeywords() []Kind {
	res := make([]Kind, keyword_end-keyword_beg)

	for i := range int(keyword_end) - int(keyword_beg) {
		res[i] = Kind(i + int(keyword_beg))
	}

	return res
}
func (t Token) IsComparisonOp() bool {
	return (comparison_op_beg < t.Kind && t.Kind < comparison_op_end) ||
		t.Kind == IS
}

func IsNotOp(r rune) bool {
	for i := operator_beg + 1; i < operator_end; i++ {
		t := tokens[i]
		if t != "" && r == rune(t[0]) {
			return false
		}
	}

	return true
}
