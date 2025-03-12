package lexer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/token"
)

func TestStatement(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple if statement",
			input: "{%if name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple switch statement",
			input: "{%switch name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.SWITCH},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple case statement",
			input: "{%case value%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.CASE},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "value"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple default statement",
			input: "{%default%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.DEFAULT},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple end statement",
			input: "{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Simple genif statement",
			input: "{%genif name%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.GENIF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "If with equal expression",
			input: "{%if name==3%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.EQL},
				{Typ: token.INT, Val: "3"},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Empty if block",
			input: "{%if name%}{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "If block with body",
			input: "{%if name%}\nText\n{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
				{Typ: token.LBR, Val: "\n"},
				{Typ: token.TEXT, Val: "Text\n"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		// {
		// 	name:  "If block with indentation",
		// 	input: "Text\n\t{%if name%}\n\tText\n\t{%end%}",
		// 	expected: []token.Token{
		// 		{Typ: token.TEXT, Val: "Text\n"},
		// 		{Typ: token.WS, Val: "\t"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.IF},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "name"},
		// 		{Typ: token.RSTMT},
		// 		{Typ: token.LBR, Val: "\n"},
		// 		{Typ: token.TEXT, Val: "\tText\n"},
		// 		{Typ: token.WS, Val: "\t"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.END},
		// 		{Typ: token.RSTMT},
		// 	},
		// },
		{
			name:  "Emtpy switch block",
			input: "{%switch name%}\n{%case a%}\n{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.SWITCH},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
				{Typ: token.LBR, Val: "\n"},
				{Typ: token.LSTMT},
				{Typ: token.CASE},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "a"},
				{Typ: token.RSTMT},
				{Typ: token.LBR, Val: "\n"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		{
			name:  "Switch block",
			input: "{%switch name%}\n{%case a%}\nText\n{%end%}",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.SWITCH},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
				{Typ: token.RSTMT},
				{Typ: token.LBR, Val: "\n"},
				{Typ: token.LSTMT},
				{Typ: token.CASE},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "a"},
				{Typ: token.RSTMT},
				{Typ: token.LBR, Val: "\n"},
				{Typ: token.TEXT, Val: "Text\n"},
				{Typ: token.LSTMT},
				{Typ: token.END},
				{Typ: token.RSTMT},
			},
		},
		// {
		// 	name:  "Switch block with indentation",
		// 	input: "Text\n\t{%switch name%}\n\t{%case a%}\n\tText\n\t{%end%}",
		// 	expected: []token.Token{
		// 		{Typ: token.TEXT, Val: "Text\n"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.SWITCH},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "name"},
		// 		{Typ: token.RSTMT},
		// 		{Typ: token.LBR, Val: "\n"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.CASE},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "a"},
		// 		{Typ: token.RSTMT},
		// 		{Typ: token.LBR, Val: "\n"},
		// 		{Typ: token.TEXT, Val: "Text\n"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.END},
		// 		{Typ: token.RSTMT},
		// 	},
		// },
	}
	runTestCases(t, testCases)
}

func TestStatementEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Unclosed statement",
			input: "{%",
			expected: []token.Token{
				{Typ: token.LSTMT},
			},
		},
		{
			name:  "Unclosed if statement",
			input: "{%if",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
			},
		},
		{
			name:  "Text after unclosed if statement",
			input: "{%if",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
			},
		},
		{
			name:  "Unclosed if statement with expression",
			input: "{%if name",
			expected: []token.Token{
				{Typ: token.LSTMT},
				{Typ: token.IF},
				{Typ: token.WS, Val: " "},
				{Typ: token.IDENT, Val: "name"},
			},
		},
	}
	runTestCases(t, testCases)
}
