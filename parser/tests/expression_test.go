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
				&parser.TextNode{
					Val: []string{"hello"},
				},
			},
		},
		{
			name:  "Single expression with var",
			input: "{{x}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.Ident{Name: "x"},
				},
			},
		},
		{
			name:  "Int literal",
			input: "{{1}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.NumberLit{
						Value: value.NumberValue(1),
					},
				},
			},
		},
		{
			name:  "Negative int literal",
			input: "{{-1}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.NumberLit{
						Value: value.NumberValue(-1),
					},
				},
			},
		},
		{
			name:  "Float literal",
			input: "{{1.1}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.NumberLit{
						Value: value.NumberValue(1.1),
					},
				},
			},
		},
		{
			name:  "Whitespaces with var",
			input: "{{ x }}",
			expected: []parser.Node{
				&parser.ExprNode{
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
				&parser.TextNode{
					Val: []string{"Hello, "},
				},
				&parser.ExprNode{
					Body: &parser.Ident{Name: "name"},
				},
				&parser.TextNode{
					Val: []string{"\n"},
				},
				&parser.ExprNode{
					Body: &parser.Ident{Name: "var"},
				},
				&parser.TextNode{
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
				&parser.ExprNode{
					Body: &parser.ParenExpr{
						Expr: &parser.Ident{
							Name: "age",
						},
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

func TestExpressionsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Empty expression block",
			input:    "{{}}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Whitespaces inside expression block",
			input:    "{{ \t }}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression interrupted with EOF",
			input:    "{{",
			expected: []parser.Node{},
			// errExpected: parser.ExpectedTokensError{
			// 	Tokens: []token.Kind{token.REXPR},
			// },
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression with var interrupted with EOF",
			input:    "{{var",
			expected: []parser.Node{},
			errExpected: parser.ExpectedTokensError{
				Tokens: []token.Kind{token.REXPR},
			},
			// errExpected: parser.Error{
			// 	Typ: parser.ErrExpressionExpected,
			// },
		},
		{
			name:     "Expression interrupted with line break",
			input:    "{{\n",
			expected: []parser.Node{},
			// errExpected: parser.ExpectedTokensError{
			// 	Tokens: []token.Kind{token.REXPR},
			// },
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression interrupted with line break with text and expression after",
			input:    "{{\nText{{1}}",
			expected: []parser.Node{},
			// errExpected: parser.ExpectedTokensError{
			// 	Tokens: []token.Kind{token.REXPR},
			// },
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression with open statement",
			input:    "{{ {% }}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression with statement",
			input:    "{{ {%%} }}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression with closing statement",
			input:    "{{ %} }}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Expression with closing comment",
			input:    "{{ #} }}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
		{
			name:     "Keyword 'and' as ident name",
			input:    "{{and}}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrExpressionExpected,
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperators(t *testing.T) {
	testCases := []testCase{
		{
			name:  "String literals equal",
			input: `{{"a"=="b"}}`,
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.StringLit{
							Quote: '"',
							Value: value.StringValue("a"),
						},
						Op: parser.Kw{
							Kind: token.EQL,
						},
						Y: &parser.StringLit{
							Quote: '"',
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.EQL,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.IS,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.ISNOT,
						},
						Y: &parser.Ident{
							Name: "b",
						},
					},
				},
			},
		},
		{
			name:  "Vars not equal (!=)",
			input: "{{a!=b}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.NEQL,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.UnaryExpr{
						Op: parser.Kw{
							Kind: token.NOT,
						},
						Expr: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.UnaryExpr{
						Op: parser.Kw{
							Kind: token.EXCL,
						},
						Expr: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.LAND,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.LOR,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.AND,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.OR,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Ident{
								Name: "a",
							},
							Op: parser.Kw{
								Kind: token.LOR,
							},
							Y: &parser.Ident{
								Name: "b",
							},
						},
						Op: parser.Kw{
							Kind: token.LOR,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.BinaryExpr{
							X: &parser.Ident{
								Name: "a",
							},
							Op: parser.Kw{
								Kind: token.LAND,
							},
							Y: &parser.Ident{
								Name: "b",
							},
						},
						Op: parser.Kw{
							Kind: token.LAND,
						},
						Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "a",
						},
						Op: parser.Kw{
							Kind: token.OR,
						},
						Y: &parser.BinaryExpr{
							X: &parser.Ident{
								Name: "b",
							},
							Op: parser.Kw{
								Kind: token.AND,
							},
							Y: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.ParenExpr{
							Expr: &parser.BinaryExpr{
								X: &parser.Ident{Name: "a"},
								Op: parser.Kw{
									Kind: token.OR,
								},
								Y: &parser.Ident{Name: "b"},
							},
						},
						Op: parser.Kw{
							Kind: token.AND,
						},
						Y: &parser.Ident{Name: "c"},
					},
				},
			},
		},
		{
			name:  "Parens changing precedence 2",
			input: "{{a and (b or c)}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.Ident{Name: "a"},
						Op: parser.Kw{
							Kind: token.AND,
						},
						Y: &parser.ParenExpr{
							Expr: &parser.BinaryExpr{
								X: &parser.Ident{Name: "b"},
								Op: parser.Kw{
									Kind: token.OR,
								},
								Y: &parser.Ident{Name: "c"},
							},
						},
					},
				},
			},
		},
		{
			name:  "Nested parens",
			input: "{{((a or b) and c) or d}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.BinaryExpr{
						X: &parser.ParenExpr{
							Expr: &parser.BinaryExpr{
								X: &parser.ParenExpr{
									Expr: &parser.BinaryExpr{
										X: &parser.Ident{Name: "a"},
										Op: parser.Kw{
											Kind: token.OR,
										},
										Y: &parser.Ident{Name: "b"},
									},
								},
								Op: parser.Kw{
									Kind: token.AND,
								},
								Y: &parser.Ident{Name: "c"},
							},
						},
						Op: parser.Kw{
							Kind: token.OR,
						},
						Y: &parser.Ident{Name: "d"},
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
			// FIXME: not valid
			name:  "Simple ternary",
			input: `{{flag?1:2}}`,
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.Ident{
							Name: "flag",
						},
						Do: parser.Kw{
							Kind: token.QUESTION,
						},
						TrueExpr: &parser.NumberLit{
							Value: value.NumberValue(1),
						},
						Else: parser.Kw{
							Kind: token.COLON,
						},
						FalseExpr: &parser.NumberLit{
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
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.Ident{
							Name: "flag",
						},
						Do: parser.Kw{
							Kind: token.DO,
						},
						TrueExpr: &parser.Ident{
							Name: "a",
						},
						Else: parser.Kw{
							Kind: token.ELSE,
						},
						FalseExpr: &parser.Ident{
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
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.BinaryExpr{
							X: &parser.Ident{
								Name: "flag",
							},
							Op: parser.Kw{
								Kind: token.EQL,
							},
							Y: &parser.NumberLit{
								Value: value.NumberValue(3),
							},
						},
						Do: parser.Kw{
							Kind: token.QUESTION,
						},
						TrueExpr: &parser.NumberLit{
							Value: value.NumberValue(1),
						},
						Else: parser.Kw{
							Kind: token.COLON,
						},
						FalseExpr: &parser.NumberLit{
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
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.Ident{
							Name: "flag",
						},
						Do: parser.Kw{
							Kind: token.QUESTION,
						},
						TrueExpr: &parser.NumberLit{
							Value: value.NumberValue(1),
						},
						Else: parser.Kw{
							Kind: token.COLON,
						},
						FalseExpr: &parser.NumberLit{
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
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.Ident{
							Name: "flag",
						},
						Do: parser.Kw{
							Kind: token.QUESTION,
						},
						TrueExpr: &parser.Ident{
							Name: "name",
						},
						Else: parser.Kw{
							Kind: token.COLON,
						},
						FalseExpr: &parser.Ident{
							Name: "age",
						},
					},
				},
			},
		},
		{
			name:  "Nested ternaries",
			input: `{{ flag ? bar ? 1 : 3 : 2 }}`,
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.TernaryExpr{
						Condition: &parser.Ident{
							Name: "flag",
						},
						Do: parser.Kw{
							Kind: token.QUESTION,
						},
						TrueExpr: &parser.TernaryExpr{
							Condition: &parser.Ident{
								Name: "bar",
							},
							Do: parser.Kw{
								Kind: token.QUESTION,
							},
							TrueExpr: &parser.NumberLit{
								Value: value.NumberValue(1),
							},
							Else: parser.Kw{
								Kind: token.COLON,
							},
							FalseExpr: &parser.NumberLit{
								Value: value.NumberValue(3),
							},
						},
						Else: parser.Kw{
							Kind: token.COLON,
						},
						FalseExpr: &parser.NumberLit{
							Value: value.NumberValue(2),
						},
					},
				},
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
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.FilterExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
						Filter: parser.Ident{
							Name: "upper",
						},
					},
				},
			},
		},
		{
			name:  "Filter with whitespaces",
			input: "{{ name -> upper }}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.FilterExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
						Filter: parser.Ident{
							Name: "upper",
						},
					},
				},
			},
		},
		{
			name:  "Nested filters",
			input: "{{name->upper->camel}}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.FilterExpr{
						Expr: &parser.FilterExpr{
							Expr: &parser.Ident{
								Name: "name",
							},
							Filter: parser.Ident{
								Name: "upper",
							},
						},
						Filter: parser.Ident{
							Name: "camel",
						},
					},
				},
			},
		},
		{
			name:  "Nested filters with whitespaces",
			input: "{{ name  -> upper -> camel }}",
			expected: []parser.Node{
				&parser.ExprNode{
					Body: &parser.FilterExpr{
						Expr: &parser.FilterExpr{
							Expr: &parser.Ident{
								Name: "name",
							},
							Filter: parser.Ident{
								Name: "upper",
							},
						},
						Filter: parser.Ident{
							Name: "camel",
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}
