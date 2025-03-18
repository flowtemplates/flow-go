package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
)

func TestIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Truthy if statement",
			str:      "{%if var%}\ntext\n{%end%}",
			expected: "\ntext\n",
			input: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{"\n", "text", "\n"},
						},
					},
					Else: nil,
				},
			},
			scope: renderer.Scope{
				"var": true,
			},
			errExpected: false,
		},
		{
			name:     "Falsy if statement",
			str:      "{%if var%}\ntext\n{%end%}",
			expected: "",
			input: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{"\n", "text", "\n"},
						},
					},
					Else: nil,
				},
			},
			scope: renderer.Scope{
				"var": false,
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
