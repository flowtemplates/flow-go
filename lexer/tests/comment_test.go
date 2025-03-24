package lexer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/token"
)

func TestComments(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty comment",
			input: `{##}`,
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.RCOMM},
			},
		},
		{
			name:  "Single comment",
			input: "{#no comments..#}",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: `no comments..`},
				{Kind: token.RCOMM},
			},
		},
		{
			name:  "Single comment with whitespaces",
			input: "{# no comments.. #}",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: ` no comments.. `},
				{Kind: token.RCOMM},
			},
		},
		{
			name: "Multiline comment",
			input: `
{# line 1
line 2

line 3 #}`[1:],
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: " line 1\nline 2\n\nline 3 "},
				{Kind: token.RCOMM},
			},
		},
		{
			name: "Multiple comments",
			input: `
{# line 1 #}
{# line 2 #}
{# line 3 #}`[1:],
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: " line 1 "},
				{Kind: token.RCOMM},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: " line 2 "},
				{Kind: token.RCOMM},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: " line 3 "},
				{Kind: token.RCOMM},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestCommentsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Only comm begin",
			input: "{#",
			expected: []token.Token{
				{Kind: token.LCOMM},
			},
		},
		{
			name:  "Comment interrupted with EOF",
			input: "{# Some content",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: " Some content"},
			},
		},
		{
			name:  "Comment with expression inside",
			input: "{#{{}}#}",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: "{{}}"},
				{Kind: token.RCOMM},
			},
		},
		{
			name:  "Comment with statement inside",
			input: "{#{%%}#}",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: "{%%}"},
				{Kind: token.RCOMM},
			},
		},
		{
			name:  "Unclosed comment with statement inside",
			input: "{#{%%}",
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: "{%%}"},
			},
		},
		{
			name: "Unclosed multiline comment with statement inside",
			input: `
{#
{%if%}
`[1:],
			expected: []token.Token{
				{Kind: token.LCOMM},
				{Kind: token.COMM_TEXT, Val: "\n{%if%}\n"},
			},
		},
	}
	runTestCases(t, testCases)
}
