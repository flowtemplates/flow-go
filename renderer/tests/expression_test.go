package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
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
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Int literal",
			str:      "{{1}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Value: value.NumberValue(1),
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Float literal",
			str:      "{{1.1}}",
			expected: "1.1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Value: value.NumberValue(1.1),
					},
				},
			},
			scope:       renderer.Scope{},
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
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "String literal",
			str:      `{{"word"}}`,
			expected: "word",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Lit{
						Value: value.StringValue("word"),
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Addition",
			str:      "{{123+2}}",
			expected: "125",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Lit{
							Value: value.NumberValue(123),
						},
						Op: token.ADD,
						Y: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		// {
		// 	name:     "Subtraction",
		// 	str:      "{{123-10}}",
		// 	expected: "113",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.SUB,
		// 				Y: parser.Lit{
		// 					Value: value.NumberValue(10),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	scope:       renderer.Scope{},
		// 	errExpected: false,
		// },
		{
			name:     "Expression with string var",
			str:      "{{name}}",
			expected: "useuse",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
			},
			scope: renderer.Scope{
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
			scope: renderer.Scope{
				"name": "world",
				"flow": "flow",
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}

func TestTernaryExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Simple ternary with true",
			str:      "{{true?1:2}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "true",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Simple ternary with false",
			str:      "{{false?1:2}}",
			expected: "2",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "false",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Simple ternary with true and some text around",
			str:      "arr[{{false?1:2}}]",
			expected: "arr[2]",
			input: []parser.Node{
				parser.Text{
					Val: []string{"arr["},
				},
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "false",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
				parser.Text{
					Val: []string{"]"},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Simple ternary",
			str:      "{{flag?1:2}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope: renderer.Scope{
				"flag": true,
			},
			errExpected: false,
		},
		{
			name:     "Do-else ternary",
			str:      "{{flag do 1 else 2}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.DO,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.ELSE,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope: renderer.Scope{
				"flag": true,
			},
			errExpected: false,
		},
		{
			name:     "Ternary with truthy number condition",
			str:      "{{1?1:2}}",
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Lit{
							Value: value.NumberValue(1),
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Ternary with falsy number condition",
			str:      "{{0?1:2}}",
			expected: "2",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Lit{
							Value: value.NumberValue(0),
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Ternary with truthy string condition",
			str:      `{{"a"?1:2}}`,
			expected: "1",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Lit{
							Value: value.StringValue("a"),
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Ternary with falsy string condition",
			str:      `{{""?1:2}}`,
			expected: "2",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Lit{
							Value: value.StringValue(""),
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.NumberValue(1),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Ternary with 3 vars",
			str:      "{{flag?a:b}}",
			expected: "foo",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Ident{
							Name: "a",
						},
						Else: token.COLON,
						FalseExpr: parser.Ident{
							Name: "b",
						},
					},
				},
			},
			scope: renderer.Scope{
				"flag": true,
				"a":    "foo",
				"b":    "bar",
			},
			errExpected: false,
		},
		{
			name:     "Ternary with truthy equal",
			str:      `{{flag + 1==3?"foo":"bar"}}`,
			expected: "foo",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.BinaryExpr{
							X: parser.BinaryExpr{
								X: parser.Ident{
									Name: "flag",
								},
								Op: token.ADD,
								Y: parser.Lit{
									Value: value.NumberValue(1),
								},
							},
							Op: token.EQL,
							Y: parser.Lit{
								Value: value.NumberValue(3),
							},
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.StringValue("foo"),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.StringValue("bar"),
						},
					},
				},
			},
			scope: renderer.Scope{
				"flag": 2,
			},
			errExpected: false,
		},
		{
			name:     "Ternary with falsy equal",
			str:      `{{flag+1==3?"foo":"bar"}}`,
			expected: "bar",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.BinaryExpr{
							X: parser.BinaryExpr{
								X: parser.Ident{
									Name: "flag",
								},
								Op: token.ADD,
								Y: parser.Lit{
									Value: value.NumberValue(1),
								},
							},
							Op: token.EQL,
							Y: parser.Lit{
								Value: value.NumberValue(3),
							},
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Value: value.StringValue("foo"),
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Value: value.StringValue("bar"),
						},
					},
				},
			},
			scope: renderer.Scope{
				"flag": 20,
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
