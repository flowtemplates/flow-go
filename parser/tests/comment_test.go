package parser_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
)

func TestComments(t *testing.T) {
	testCases := []testCase{
		{
			name: "Empty comment",
			input: `
{##}`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val: "",
				},
			},
		},
		{
			name: "Simple single line comment",
			input: `
{# asd #}`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val: "asd",
				},
			},
		},
		{
			name: "Single line comment",
			input: `
{# asd #}`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val: "asd",
				},
			},
		},
		{
			name: "Single line comment with line break after",
			input: `
{# asd #}
`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val:    "asd",
					PostLB: "\n",
				},
			},
		},
		{
			name: "Single line comment with whitespaces after",
			input: `
{# asd #}  `[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val:    "asd",
					PostLB: "",
				},
			},
		},
		{
			name: "Single line comment with whitespaces and line break after",
			input: `
{# asd #}
`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val:    "asd",
					PostLB: "\n",
				},
			},
		},
		{
			name: "Multiline comment",
			input: `
{#
123
Text
#}`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val: "\n123\nText\n",
				},
			},
		},
		{
			name: "Multiple comments",
			input: `
{# Text #}
{# 2 #}
{# 3 #}
`[1:],
			expected: []parser.Node{
				&parser.CommNode{
					Val:    "Text",
					PostLB: "\n",
				},
				&parser.CommNode{
					Val:    "2",
					PostLB: "\n",
				},
				&parser.CommNode{
					Val:    "3",
					PostLB: "\n",
				},
			},
		},
	}
	runTestCases(t, testCases)
}
