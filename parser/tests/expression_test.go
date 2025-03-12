package parser_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			str:  "hello",
			input: []token.Token{
				{Kind: token.TEXT, Val: "hello"},
			},
			expected: []parser.Node{
				parser.Text{
					Val: []string{"hello"},
				},
			},
		},
		{
			name: "Single expression with var",
			str:  "{{x}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "x"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{Name: "x"},
				},
			},
		},
		{
			name: "Whitespaces with var",
			str:  "{{ x }}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "x"},
				{Kind: token.WS, Val: " "},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{
						Name: "x",
					},
				},
			},
		},
		{
			name: "Expressions between text",
			str:  "Hello, {{name}}\n{{var }}Text",
			input: []token.Token{
				{Kind: token.TEXT, Val: "Hello, "},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.REXPR},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: " "},
				{Kind: token.REXPR},
				{Kind: token.TEXT, Val: "Text"},
			},
			expected: []parser.Node{
				parser.Text{
					Val: []string{"Hello, "},
				},
				parser.ExprBlock{
					Body: &parser.Ident{Name: "name"},
				},
				parser.Text{
					Val: []string{"\n"},
				},
				parser.ExprBlock{
					Body: &parser.Ident{Name: "var"},
				},
				parser.Text{
					Val: []string{"Text"},
				},
			},
		},
		{
			name: "Addition",
			str:  "{{123+age}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "123"},
				{Kind: token.ADD},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.ADD,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name: "Addition with whitespaces",
			str:  "{{123 + age}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "123"},
				{Kind: token.WS, Val: " "},
				{Kind: token.ADD},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.ADD,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name: "Subtraction",
			str:  "{{123-age}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "123"},
				{Kind: token.SUB},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.SUB,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name: "Multiply",
			str:  "{{123*age}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "123"},
				{Kind: token.MUL},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.MUL,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},

		{
			name: "Multiply",
			str:  "{{123*age}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.INT, Val: "123"},
				{Kind: token.MUL},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Lit{
							Val: "123",
							Typ: token.INT,
						},
						Op: token.MUL,
						Y: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name: "Division",
			str:  "{{age/2}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.DIV},
				{Kind: token.INT, Val: "2"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "age",
						},
						Op: token.DIV,
						Y: &parser.Lit{
							Val: "2",
							Typ: token.INT,
						},
					},
				},
			},
		},
		{
			name: "Redundant parens",
			str:  "{{(age)}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.LPAREN},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.RPAREN},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{
						Name: "age",
					},
				},
			},
		},
		{
			name: "Parens * int",
			str:  "{{(1+2)*3}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.LPAREN},
				{Kind: token.INT, Val: "1"},
				{Kind: token.ADD},
				{Kind: token.INT, Val: "2"},
				{Kind: token.RPAREN},
				{Kind: token.MUL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Lit{
								Val: "1",
								Typ: token.INT,
							},
							Op: token.ADD,
							Y: &parser.Lit{
								Val: "2",
								Typ: token.INT,
							},
						},
						Op: token.MUL,
						Y: &parser.Lit{
							Val: "3",
							Typ: token.INT,
						},
					},
				},
			},
		},
		{
			name: "Parens * int (with whitespaces)",
			str:  "{{(1 + 2) * 3}}",
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.LPAREN},
				{Kind: token.INT, Val: "1"},
				{Kind: token.WS, Val: " "},
				{Kind: token.ADD},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "2"},
				{Kind: token.RPAREN},
				{Kind: token.WS, Val: " "},
				{Kind: token.MUL},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.BinaryExpr{
							X: parser.Lit{
								Val: "1",
								Typ: token.INT,
							},
							Op: token.ADD,
							Y: parser.Lit{
								Val: "2",
								Typ: token.INT,
							},
						},
						Op: token.MUL,
						Y: parser.Lit{
							Val: "3",
							Typ: token.INT,
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestTernaries(t *testing.T) {
	testCases := []testCase{
		{
			name: "Simple ternary",
			str:  `{{flag?1:2}}`,
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.QUESTION},
				{Kind: token.INT, Val: "1"},
				{Kind: token.COLON},
				{Kind: token.INT, Val: "2"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Typ: token.INT,
							Val: "1",
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Typ: token.INT,
							Val: "2",
						},
					},
				},
			},
		},
		{
			name: "Simple do-else ternary",
			str:  `{{flag do a else b}}`,
			input: []token.Token{
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
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.DO,
						TrueExpr: parser.Ident{
							Name: "a",
						},
						Else: token.ELSE,
						FalseExpr: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name: "Ternary with equality condition",
			str:  `{{flag==3?1:2}}`,
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.EQL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.QUESTION},
				{Kind: token.INT, Val: "1"},
				{Kind: token.COLON},
				{Kind: token.INT, Val: "2"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.BinaryExpr{
							X: parser.Ident{
								Name: "flag",
							},
							Op: token.EQL,
							Y: parser.Lit{
								Typ: token.INT,
								Val: "3",
							},
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Typ: token.INT,
							Val: "1",
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Typ: token.INT,
							Val: "2",
						},
					},
				},
			},
		},
		{
			name: "Ternary with whitespaces",
			str:  `{{ flag ? 1 : 2 }}`,
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.WS, Val: " "},
				{Kind: token.QUESTION},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "1"},
				{Kind: token.WS, Val: " "},
				{Kind: token.COLON},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "2"},
				{Kind: token.WS, Val: " "},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Lit{
							Typ: token.INT,
							Val: "1",
						},
						Else: token.COLON,
						FalseExpr: parser.Lit{
							Typ: token.INT,
							Val: "2",
						},
					},
				},
			},
		},
		{
			name: "Ternary with vars",
			str:  `{{flag?name:age}}`,
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.QUESTION},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.COLON},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.Ident{
							Name: "name",
						},
						Else: token.COLON,
						FalseExpr: parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name: "Ternary with complex expressions",
			str:  `{{flag?(1/2)+name:age*3}}`,
			input: []token.Token{
				{Kind: token.LEXPR},
				{Kind: token.IDENT, Val: "flag"},
				{Kind: token.QUESTION},
				{Kind: token.LPAREN},
				{Kind: token.INT, Val: "1"},
				{Kind: token.DIV},
				{Kind: token.INT, Val: "2"},
				{Kind: token.RPAREN},
				{Kind: token.ADD},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.COLON},
				{Kind: token.IDENT, Val: "age"},
				{Kind: token.MUL},
				{Kind: token.INT, Val: "3"},
				{Kind: token.REXPR},
			},
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.Ident{
							Name: "flag",
						},
						Do: token.QUESTION,
						TrueExpr: parser.BinaryExpr{
							X: parser.BinaryExpr{
								X: parser.Lit{
									Typ: token.INT,
									Val: "1",
								},
								Op: token.DIV,
								Y: parser.Lit{
									Typ: token.INT,
									Val: "2",
								},
							},
							Op: token.ADD,
							Y: parser.Ident{
								Name: "name",
							},
						},
						Else: token.COLON,
						FalseExpr: parser.BinaryExpr{
							X: parser.Ident{
								Name: "age",
							},
							Op: token.MUL,
							Y: parser.Lit{
								Typ: token.INT,
								Val: "3",
							},
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}
