package parser_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Plain text",
			input: "hello",
			expected: []parser.Node{
				parser.Text{
					Val: []string{"hello"},
				},
			},
		},
		{
			name:  "Single expression with var",
			input: "{{x}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{Name: "x"},
				},
			},
		},
		{
			name:  "Whitespaces with var",
			input: "{{ x }}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{
						Name: "x",
					},
				},
			},
		},
		{
			name:  "Expressions between text",
			input: "Hello, {{name}}\n{{var }}Text",
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
		// {
		// 	name: "Addition",
		// 	str:  "{{123+age}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.INT, Val: "123"},
		// 		{Kind: token.ADD},
		// 		{Kind: token.IDENT, Val: "age"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.ADD,
		// 				Y: &parser.Ident{
		// 					Name: "age",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Addition with whitespaces",
		// 	str:  "{{123 + age}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.INT, Val: "123"},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.ADD},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.IDENT, Val: "age"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.ADD,
		// 				Y: &parser.Ident{
		// 					Name: "age",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Subtraction",
		// 	str:  "{{123-age}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.INT, Val: "123"},
		// 		{Kind: token.SUB},
		// 		{Kind: token.IDENT, Val: "age"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.SUB,
		// 				Y: &parser.Ident{
		// 					Name: "age",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Multiply",
		// 	str:  "{{123*age}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.INT, Val: "123"},
		// 		{Kind: token.MUL},
		// 		{Kind: token.IDENT, Val: "age"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.MUL,
		// 				Y: &parser.Ident{
		// 					Name: "age",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Division",
		// 	str:  "{{age/2}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.IDENT, Val: "age"},
		// 		{Kind: token.DIV},
		// 		{Kind: token.INT, Val: "2"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.Ident{
		// 					Name: "age",
		// 				},
		// 				Op: token.DIV,
		// 				Y: &parser.Lit{
		// 					Value: value.NumberValue(2),
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			name:  "Redundant parens",
			input: "{{(age)}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: &parser.Ident{
						Name: "age",
					},
				},
			},
		},
		// {
		// 	name: "Parens * int",
		// 	str:  "{{(1+2)*3}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.LPAREN},
		// 		{Kind: token.INT, Val: "1"},
		// 		{Kind: token.ADD},
		// 		{Kind: token.INT, Val: "2"},
		// 		{Kind: token.RPAREN},
		// 		{Kind: token.MUL},
		// 		{Kind: token.INT, Val: "3"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: &parser.BinaryExpr{
		// 				X: &parser.BinaryExpr{
		// 					X: &parser.Lit{
		// 						Value: value.NumberValue(1),
		// 					},
		// 					Op: token.ADD,
		// 					Y: &parser.Lit{
		// 						Value: value.NumberValue(2),
		// 					},
		// 				},
		// 				Op: token.MUL,
		// 				Y: &parser.Lit{
		// 					Value: value.NumberValue(3),
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "Parens * int (with whitespaces)",
		// 	str:  "{{(1 + 2) * 3}}",
		// 	input: []token.Token{
		// 		{Kind: token.LEXPR},
		// 		{Kind: token.LPAREN},
		// 		{Kind: token.INT, Val: "1"},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.ADD},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.INT, Val: "2"},
		// 		{Kind: token.RPAREN},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.MUL},
		// 		{Kind: token.WS, Val: " "},
		// 		{Kind: token.INT, Val: "3"},
		// 		{Kind: token.REXPR},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.BinaryExpr{
		// 					X: parser.Lit{
		// 						Value: value.NumberValue(1),
		// 					},
		// 					Op: token.ADD,
		// 					Y: parser.Lit{
		// 						Value: value.NumberValue(2),
		// 					},
		// 				},
		// 				Op: token.MUL,
		// 				Y: parser.Lit{
		// 					Value: value.NumberValue(3),
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
	runTestCases(t, testCases)
}

func TestOperators(t *testing.T) {
	testCases := []testCase{
		{
			name:  "String literals equal",
			input: `{{"a"=="b"}}`,
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Lit{
							Value: value.StringValue("a"),
						},
						Op: token.EQL,
						Y: parser.Lit{
							Value: value.StringValue("b"),
						},
					},
				},
			},
		},
		{
			name:  "Vars equal (==)",
			input: "{{a==b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.EQL,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars is",
			input: "{{a is b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.IS,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars is not",
			input: "{{a is not b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.IS,
						Y: parser.UnaryExpr{
							Op: token.NOT,
							X: parser.Ident{
								Name: "b",
							},
						},
					},
				},
			},
		},
		{
			name:  "Vars not equal (!=)",
			input: "{{a!=b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.NEQL,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars not",
			input: "{{not flag}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.UnaryExpr{
						Op: token.NOT,
						X: parser.Ident{
							Name: "flag",
						},
					},
				},
			},
		},
		{
			name:  "Vars not (!)",
			input: "{{!flag}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.UnaryExpr{
						Op: token.EXCL,
						X: parser.Ident{
							Name: "flag",
						},
					},
				},
			},
		},
		{
			name:  "Vars and (&&)",
			input: "{{a&&b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.LAND,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars or (||)",
			input: "{{a||b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.LOR,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars and",
			input: "{{a and b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.AND,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars or",
			input: "{{a or b}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.OR,
						Y: parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Multiple ||",
			input: "{{a||b||c}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.BinaryExpr{
							X: parser.Ident{
								Name: "a",
							},
							Op: token.LOR,
							Y: parser.Ident{
								Name: "b",
							},
						},
						Op: token.LOR,
						Y: parser.Ident{
							Name: "c",
						},
					},
				},
			},
		},
		{
			name:  "Multiple &&",
			input: "{{a&&b&&c}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.BinaryExpr{
							X: parser.Ident{
								Name: "a",
							},
							Op: token.LAND,
							Y: parser.Ident{
								Name: "b",
							},
						},
						Op: token.LAND,
						Y: parser.Ident{
							Name: "c",
						},
					},
				},
			},
		},
		{
			name:  "And precedence over or",
			input: "{{a or b and c}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "a",
						},
						Op: token.OR,
						Y: parser.BinaryExpr{
							X: parser.Ident{
								Name: "b",
							},
							Op: token.AND,
							Y: parser.Ident{
								Name: "c",
							},
						},
					},
				},
			},
		},
		{
			name:  "Parens changing precedence 1",
			input: "{{(a or b) and c}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.BinaryExpr{
							X:  parser.Ident{Name: "a"},
							Op: token.OR,
							Y:  parser.Ident{Name: "b"},
						},
						Op: token.AND,
						Y:  parser.Ident{Name: "c"},
					},
				},
			},
		},
		{
			name:  "Parens changing precedence 2",
			input: "{{a and (b or c)}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X:  parser.Ident{Name: "a"},
						Op: token.AND,
						Y: parser.BinaryExpr{
							X:  parser.Ident{Name: "b"},
							Op: token.OR,
							Y:  parser.Ident{Name: "c"},
						},
					},
				},
			},
		},
		{
			name:  "Nested parens",
			input: "{{((a or b) and c) or d}}",
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.BinaryExpr{
							X: parser.BinaryExpr{
								X:  parser.Ident{Name: "a"},
								Op: token.OR,
								Y:  parser.Ident{Name: "b"},
							},
							Op: token.AND,
							Y:  parser.Ident{Name: "c"},
						},
						Op: token.OR,
						Y:  parser.Ident{Name: "d"},
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
			name:  "Simple ternary",
			input: `{{flag?1:2}}`,
			expected: []parser.Node{
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
		},
		{
			name:  "Simple do-else ternary",
			input: `{{flag do a else b}}`,
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
			name:  "Ternary with equality condition",
			input: `{{flag==3?1:2}}`,
			expected: []parser.Node{
				parser.ExprBlock{
					Body: parser.TernaryExpr{
						Condition: parser.BinaryExpr{
							X: parser.Ident{
								Name: "flag",
							},
							Op: token.EQL,
							Y: parser.Lit{
								Value: value.NumberValue(3),
							},
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
		},
		{
			name:  "Ternary with whitespaces",
			input: `{{ flag ? 1 : 2 }}`,
			expected: []parser.Node{
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
		},
		{
			name:  "Ternary with vars",
			input: `{{flag?name:age}}`,
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
	}
	runTestCases(t, testCases)
}
