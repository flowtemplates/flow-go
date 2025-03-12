package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
	"github.com/flowtemplates/flow-go/token"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Plain text",
			str:      "Hello world",
			expected: "Hello world",
			input: []parser.Node{
				parser.Text{
					Val: []string{"Hello world"},
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Int literal",
			str:      "{{1}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Typ: token.INT,
						Val: "1",
					},
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Float literal",
			str:      "{{1.1}}",
			expected: "1.1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Typ: token.FLOAT,
						Val: "1.1",
					},
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Boolean literal",
			str:      "{{true}}",
			expected: "",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Ident{
						Name: "true",
					},
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "String literal",
			str:      `{{"word"}}`,
			expected: "word",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Typ: token.STR,
						Val: "word",
					},
				},
			},
			context:     renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Expression with string var",
			str:      "{{name}}",
			expected: "useuse",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
			},
			context: renderer.Scope{
				"name": "useuse",
			},
			errExpected: false,
		},
		{
			name:     "Multiple expressions",
			str:      "Hello {{name}}!\nFrom {{ flow }} templates",
			expected: "Hello world!\nFrom flow templates",
			input: []parser.Node{
				parser.Text{
					Val: []string{"Hello "},
				},
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
				parser.Text{
					Val: []string{"!", "\n", "From "},
				},
				parser.ExprBlock{
					Body: parser.Ident{Name: "flow"},
				},
				parser.Text{
					Val: []string{" templates"},
				},
			},
			context: renderer.Scope{
				"name": "world",
				"flow": "flow",
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}

// func TestTernaryExpressions(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:     "Simple ternary",
// 			str:      "{{true:1:2}}",
// 			expected: "1",
// 			input: []parser.Node{
// 				parser.ExprBlock{
// 					Body: parser.Ident{Name: "name"},
// 				},
// 			},
// 			context:     renderer.Scope{},
// 			errExpected: false,
// 		},
// 	}
// 	runTestCases(t, testCases)
// }
