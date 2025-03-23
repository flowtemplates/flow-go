package lexer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/token"
)

func TestIfStatement(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty condition if statement",
			input: "{%if%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple if statement",
			input: "{%if name%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple switch statement",
			input: "{%switch name%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple case statement",
			input: "{%case value%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.CASE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "value"},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple default statement",
			input: "{%default%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.DEFAULT},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple end statement",
			input: "{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Simple genif statement",
			input: "{%genif name%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.GENIF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "If with equal expression",
			input: "{%if name==3%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.EQL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "Empty if block",
			input: "{%if name%}{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "If block with body",
			input: `
{%if name%}
Text
{%end%}
`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
			},
		},
		{
			name: "If block with whitespaces after tag",
			input: `
{%if name%}  
Text
{%end%}
`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.WS, Val: "  "},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
			},
		}, {
			name: "If block with indentation",
			input: `
Text
	{%if name%}
	Text
	{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "If block with text in front",
			input: `
Text{%if name%}
Text
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "If-else block",
			input: "{%if name%}{%else%}{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "If-elseif block",
			input: "{%if name%}{%else if%}{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IF},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name:  "If-elseif-else block",
			input: "{%if name%}{%else if%}{%else%}{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IF},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.RSTMT},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "If-elseif-else block (with text between)",
			input: `
{%if name%}
Text
{%else if%}
2
{%else%}
Sometext
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IF},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "2"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.ELSE},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Sometext"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestSwitchStatement(t *testing.T) {
	testCases := []testCase{
		{
			name: "Emtpy switch block",
			input: `
{%switch name%}
{%case a%}
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.CASE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "a"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "Switch block",
			input: `
{%switch name%}
{%case a%}
Text
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.CASE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "a"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "Switch block with whitespaces after tag",
			input: `
{%switch name%}  
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.WS, Val: "  "},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "Switch block with text after tag",
			input: `
{%switch name%}asd
{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.TEXT, Val: "asd"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
		{
			name: "Switch block with indentation",
			input: `
Text
	{%switch name%}
	{%case a%}
	Text
	{%end%}`[1:],
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.LSTMT},
				{Kind: token.SWITCH},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.LSTMT},
				{Kind: token.CASE},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "a"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.WS, Val: "\t"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestStatementEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Unclosed statement",
			input: "{%",
			expected: []token.Token{
				{Kind: token.LSTMT},
			},
		},
		{
			name:  "Separe statement open",
			input: "{ %",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "{ %"},
			},
		},
		{
			name:  "Unclosed if statement",
			input: "{%if",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
			},
		},
		{
			name:  "Text after unclosed if statement",
			input: "{%if",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
			},
		},
		{
			name:  "Unclosed if statement with expression",
			input: "{%if name",
			expected: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
			},
		},
	}
	runTestCases(t, testCases)
}
