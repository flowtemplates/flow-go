package lexer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/token"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Empty input",
			input:    "",
			expected: []token.Token{},
		},
		{
			name:  "Plain text",
			input: "Hello, world!",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Hello, world!"},
			},
		},
		{
			name: "Multiline plain text",
			input: `
123
3
4
`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "123"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "3"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "4"},
				{Kind: token.LNBR, Val: "\n"},
			},
		},
		{
			name: "Multiline plain text with empty lines",
			input: `
123
3


4

`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "123"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "3"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "4"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LNBR, Val: "\n"},
			},
		},
		{
			name:  "Simple expression",
			input: "{{name}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Whitespaces inside expr",
			input: "{{ name		}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.WS, Val: "		"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Multiple expressions",
			input: "{{greeting}}, {{name}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "greeting"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: ", "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with underscores",
			input: "{{_user_name}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "_user_name"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with digits",
			input: "{{user123}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "user123"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with leading dollar sign",
			input: "{{$name}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "$name"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with dollar sign",
			input: "{{mirco$oft}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "mirco$oft"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with non-latin symbols",
			input: "{{こんにちは}} {{🙋}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "こんにちは"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: " "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "🙋"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Text before, after, and between expressions",
			input: "Hello, {{greeting}}, {{name}}! Welcome!",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Hello, "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "greeting"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: ", "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: "! Welcome!"},
			},
		},
		{
			name: "Multiline input with several blocks",
			input: `
Hello,
{{greeting}}
{{name}}{{surname}}
`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Hello,"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "greeting"},
				{Kind: token.REXPR},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "surname"},
				{Kind: token.REXPR},
				{Kind: token.LNBR, Val: "\n"},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestExpressionsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty expression",
			input: "{{}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with leading digit",
			input: "{{1user}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "1"},
				{Kind: token.IDENT, Val: "user"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with period",
			input: "{{us.er}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "us"},
				{Kind: token.PERIOD},
				{Kind: token.IDENT, Val: "er"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Var name with minus",
			input: "{{us-er}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "us"},
				{Kind: token.MINUS},
				{Kind: token.IDENT, Val: "er"},
				{Kind: token.REXPR},
			},
		},
		{
			name: "Line break inside expression",
			input: `
{{greeting
}}`[1:],
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "greeting"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Unclosed expression",
			input: "Hello, {{name",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Hello, "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
			},
		},
		{
			name:  "Only left expr",
			input: "{{",
			expected: []token.Token{
				{Kind: token.LEXPR},
			},
		},
		{
			name:  "Only right expr",
			input: "}}",
			expected: []token.Token{
				{Kind: token.REXPR},
			},
		},
		{
			name: "Text after unclosed expression",
			input: `
{{greeting
Some text`[1:],
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "greeting"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Some text"},
			},
		},
		{
			name:  "Single right expr between text",
			input: "Another text}}Some text",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Another text"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: "Some text"},
			},
		},
		{
			name: "Right expr after valid expression",
			input: `
{{name}}
}}Some text`[1:],
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: "Some text"},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestNumLiterals(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Integer value",
			input: "{{10}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "10"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Negative integer value",
			input: "{{-123}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.MINUS, Val: "-"},
				{Kind: token.INT, Val: "123"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Float value",
			input: "{{12.3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.FLOAT, Val: "12.3"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Negative float value",
			input: "{{-12.3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.MINUS, Val: "-"},
				{Kind: token.FLOAT, Val: "12.3"},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

// func TestOperations(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:  "Addittion",
// 			input: "{{seconds+1}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "seconds"},
// 				{Kind: token.ADD},
// 				{Kind: token.INT, Val: "1"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Subtraction",
// 			input: "{{age-123.2}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.SUB},
// 				{Kind: token.FLOAT, Val: "123.2"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Negative number subtraction",
// 			input: "{{age- -123.2}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.SUB},
// 				{Kind: token.WS, Val: " "},
// 				{Kind: token.SUB},
// 				{Kind: token.FLOAT, Val: "123.2"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Multiply",
// 			input: "{{age*30}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.MUL},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Multiply by negative number",
// 			input: "{{age*-30}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.MUL},
// 				{Kind: token.SUB},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Multiply by negative number",
// 			input: "{{age*-30}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.MUL},
// 				{Kind: token.SUB},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Division",
// 			input: "{{age/30}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.DIV},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Division by negative number",
// 			input: "{{age/-30}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.DIV},
// 				{Kind: token.SUB},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Single parens",
// 			input: "{{(12/2)+age}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.LPAREN},
// 				{Kind: token.INT, Val: "12"},
// 				{Kind: token.DIV},
// 				{Kind: token.INT, Val: "2"},
// 				{Kind: token.RPAREN},
// 				{Kind: token.ADD},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Two operations with parens",
// 			input: "{{(age/-30)+(12-2.2)}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.LPAREN},
// 				{Kind: token.IDENT, Val: "age"},
// 				{Kind: token.DIV},
// 				{Kind: token.SUB},
// 				{Kind: token.INT, Val: "30"},
// 				{Kind: token.RPAREN},
// 				{Kind: token.ADD},
// 				{Kind: token.LPAREN},
// 				{Kind: token.INT, Val: "12"},
// 				{Kind: token.SUB},
// 				{Kind: token.FLOAT, Val: "2.2"},
// 				{Kind: token.RPAREN},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 	}
// 	runTestCases(t, testCases)
// }

// func TestOperationsEdgeCases(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:  "Unclosed addition",
// 			input: "{{1+}}",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.INT, Val: "1"},
// 				{Kind: token.ADD},
// 				{Kind: token.REXPR},
// 			},
// 		},
// 		{
// 			name:  "Unclosed expression with addition",
// 			input: "{{1+",
// 			expected: []token.Token{
// 				{Kind: token.LEXPR},
// 				{Kind: token.INT, Val: "1"},
// 				{Kind: token.ADD},
// 			},
// 		},
// 	}
// 	runTestCases(t, testCases)
// }

func TestNumLiteralsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Integer value with unclosed expression",
			input: "{{10} name",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "10"},
				{Kind: token.RBRACE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
			},
		},
		{
			name:  "Multiple periods in float value",
			input: "{{-12.3.2}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.MINUS, Val: "-"},
				{Kind: token.FLOAT, Val: "12.3"},
				{Kind: token.PERIOD, Val: "."},
				{Kind: token.INT, Val: "2"},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStringLiterals(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple string literal in double quotes",
			input: `{{"double"}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: `"double"`},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Simple string literal in single quotes",
			input: `{{'single'}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: "'single'"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Empty string literal",
			input: `{{""}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: `""`},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "String literal with whitespaces",
			input: `{{"word1 word2  	word3"}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: `"word1 word2  	word3"`},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "String literal with numbers and booleans",
			input: `{{"123 false -22.0"}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: `"123 false -22.0"`},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStringLiteralsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "String interrupted with EOF",
			input: `{{"double`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.NOT_TERMINATED_STR, Val: `"double`},
			},
		},
		{
			name:  "String interrupted with line break",
			input: "{{\"double\n",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.NOT_TERMINATED_STR, Val: `"double`},
				{Kind: token.LNBR, Val: "\n"},
			},
		},
		{
			name:  "Multiple strings",
			input: `{{"123 falseasd" 'bsdbq12 )_ asd' }}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.STR, Val: `"123 falseasd"`},
				{Kind: token.WS, Val: " "},
				{Kind: token.STR, Val: `'bsdbq12 )_ asd'`},
				{Kind: token.WS, Val: " "},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestFilters(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple filter",
			input: "{{name->upper}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RARR},
				{Kind: token.IDENT, Val: "upper"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Nested filters",
			input: "{{name -> upper -> camel}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.WS, Val: " "},
				{Kind: token.RARR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "upper"},
				{Kind: token.WS, Val: " "},
				{Kind: token.RARR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "camel"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Filter empty params",
			input: "{{name -> upper()}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.WS, Val: " "},
				{Kind: token.RARR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "upper"},
				{Kind: token.LPAREN},
				{Kind: token.RPAREN},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Filter in expression",
			input: "{{name->upper=='UP'}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RARR},
				{Kind: token.IDENT, Val: "upper"},
				{Kind: token.EQL},
				{Kind: token.STR, Val: "'UP'"},
				{Kind: token.REXPR},
			},
		},
		// {
		// 	name:  "Filter with one param",
		// 	input: "{{name -> truncate(10)}}",
		// 	expected: []token.Token{
		// 		{Typ: token.LEXPR},
		// 		{Typ: token.IDENT, Val: "name"},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.RARR},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "truncate"},
		// 		{Typ: token.LPAREN},
		// 		{Typ: token.INT, Val: "10"},
		// 		{Typ: token.RPAREN},
		// 		{Typ: token.REXPR},
		// 	},
		// },
		// {
		// 	name:  "Filter with one named param",
		// 	input: "{{name -> truncate(length=10)}}",
		// 	expected: []token.Token{
		// 		{Typ: token.LEXPR},
		// 		{Typ: token.IDENT, Val: "name"},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.RARR},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "truncate"},
		// 		{Typ: token.LPAREN},
		// 		{Typ: token.INT, Val: "10"},
		// 		{Typ: token.RPAREN},
		// 		{Typ: token.REXPR},
		// 	},
		// },
	}
	runTestCases(t, testCases)
}

func TestFiltersEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty filter",
			input: "{{name->}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RARR},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperators(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Equal",
			input: "{{age==3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.EQL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Not",
			input: "{{not flag}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.NOT},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Not equal",
			input: "{{age!=3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.NEQL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Is",
			input: "{{age is 3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.WS, Val: " "},
				{Kind: token.IS},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "is not",
			input: "{{age is not 3}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.WS, Val: " "},
				{Kind: token.IS},
				{Kind: token.WS, Val: " "},
				{Kind: token.NOT},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Excl",
			input: "{{!flag}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.EXCL},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Greater",
			input: "{{var>1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.GRTR},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Less",
			input: "{{var<1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.LESS},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Less or equal",
			input: "{{var<=1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.LEQ},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Greater or equal",
			input: "{{var>=1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.GEQ},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "And",
			input: "{{var and 1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: " "},
				{Kind: token.AND},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "or",
			input: "{{var or 1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: " "},
				{Kind: token.OR},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "&&",
			input: "{{var&&1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.LAND},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "||",
			input: "{{var||1}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.LOR},
				{Kind: token.INT, Val: "1"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Logical expression with parens",
			input: "{{(var||1)&&false}}",
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.LPAREN},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.LOR},
				{Kind: token.INT, Val: "1"},
				{Kind: token.RPAREN},
				{Kind: token.LAND},
				{Kind: token.IDENT, Val: "false"},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestTernaryOps(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple ternary",
			input: `{{flag?1:2}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.QUESTION},
				{Kind: token.INT, Val: "1"},
				{Kind: token.COLON},
				{Kind: token.INT, Val: "2"},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Ternary with whitespaces",
			input: `{{ flag ? a : b }}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.WS, Val: " "},
				{Kind: token.QUESTION},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "a"},
				{Kind: token.WS, Val: " "},
				{Kind: token.COLON},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "b"},
				{Kind: token.WS, Val: " "},
				{Kind: token.REXPR},
			},
		},
		{
			name:  "Do-else ternary",
			input: `{{flag do a else b}}`,
			expected: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.WS, Val: " "},
				{Kind: token.DO},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "a"},
				{Kind: token.WS, Val: " "},
				{Kind: token.ELSE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "b"},
				{Kind: token.REXPR},
			},
		},
	}
	runTestCases(t, testCases)
}
