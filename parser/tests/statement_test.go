package parser_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/value"
)

func TestIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "Simple if statement",
			input: `
{%if var%}
text
{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"text",
								"\n",
							},
						},
					},
				},
			},
		},
		{
			name: "If statement with whitespaces after tag",
			input: `
{%if var%}  
text
{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"text",
								"\n",
							},
						},
					},
				},
			},
		},
		{
			name:  "Simple one line if statement",
			input: "{%if var%}text{%end%}",
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{"text"},
						},
					},
				},
			},
		},
		{
			name: "If statement with extra whitespaces",
			input: `
{% if var  %}
text
{% end %}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"text",
								"\n",
							},
						},
					},
				},
			},
		},
		{
			name: "If statement with indentation",
			input: `
	{%if var%}
	text
	{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "\t",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"\t",
								"text",
								"\n",
							},
						},
					},
					EndTag: parser.StmtTag{
						PreWs: "\t",
					},
				},
			},
		},
		{
			name: "Nested if statements",
			input: `
{%if var%}
1
{%if name%}
text
{%end%}
2
{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{"1", "\n"},
						},
						&parser.IfNode{
							IfTag: parser.StmtTagWithExpr{
								StmtTag: parser.StmtTag{
									PreWs: "",
								},
								Expr: &parser.Ident{
									Name: "name",
								},
							},
							MainBody: []parser.Node{
								&parser.TextNode{
									Val: []string{"text", "\n"},
								},
							},
						},
						&parser.TextNode{
							Val: []string{"2", "\n"},
						},
					},
				},
			},
		},
		{
			name: "Simple if-elseif statement",
			input: `
{%if bar%}
1
{%else if flag%}
2
{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "bar",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"1",
								"\n",
							},
						},
					},
					ElseIfs: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "flag",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"2",
										"\n",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "If-elseif-else statement",
			input: `
{%if bar%}
1
{%else if flag%}
2
{%else%}
3
{%end%}`[1:],
			expected: []parser.Node{
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "bar",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"1",
								"\n",
							},
						},
					},
					ElseIfs: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "flag",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"2",
										"\n",
									},
								},
							},
						},
					},
					ElseBody: parser.Clause{
						Body: []parser.Node{
							&parser.TextNode{
								Val: []string{
									"3",
									"\n",
								},
							},
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestIfStatementsEdgeCases(t *testing.T) {
	testCases := []testCase{
		{
			name:     "If statement without end tag",
			input:    "{%if var%}",
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrEndExpected,
			},
		},
		// {
		// 	name:     "Empty condition if statement",
		// 	input:    "{%if %}{%end%}",
		// 	expected: []parser.Node{},
		// 	errExpected: parser.Error{
		// 		Typ: parser.ErrEndExpected,
		// 	},
		// },
		{
			name: "If statement with body without end tag",
			input: `
{%if var%}
Some {{text}}`[1:],
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrEndExpected,
			},
		},
		{
			name: "If statement without end keyword unclosed",
			input: `
{%if var%}
{%`[1:],
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrKeywordExpected,
			},
		},
		{
			name: "If statement with unclosed end tag",
			input: `
{%if var%}
{% end`[1:],
			expected: []parser.Node{},
			errExpected: parser.ExpectedTokensError{
				Tokens: []token.Kind{token.RSTMT},
			},
		},
		{
			name: "If statement without end keyword unclosed with body",
			input: `
{%if var%}
{{text}}
{%`[1:],
			expected: []parser.Node{},
			errExpected: parser.Error{
				Typ: parser.ErrKeywordExpected,
			},
		},
		{
			name: "If statement with unclosed end tag with body",
			input: `
{%if var%}
{{text}}
{% end`[1:],
			expected: []parser.Node{},
			errExpected: parser.ExpectedTokensError{
				Tokens: []token.Kind{token.RSTMT},
			},
		},
		{
			name: "If statement with text in front",
			input: `
Text{%if var%}
123
{% end %}`[1:],
			expected: []parser.Node{
				&parser.TextNode{
					Val: []string{"Text"},
				},
				&parser.IfNode{
					IfTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Expr: &parser.Ident{
							Name: "var",
						},
					},
					MainBody: []parser.Node{
						&parser.TextNode{
							Val: []string{
								"123",
								"\n",
							},
						},
					},
				},
			},
			// errExpected: parser.Error{
			// 	Typ: parser.ErrUnexpectedBeforeStmt,
			// },
		},
	}
	runTestCases(t, testCases)
}

func TestGenIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple genif",
			input: "{%genif var%}",
			expected: []parser.Node{
				&parser.StmtNode{
					StmtTagWithKw: parser.StmtTagWithKw{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Kw: parser.Kw{
							Kind: token.GENIF,
						},
					},
					Expr: &parser.Ident{
						Name: "var",
					},
				},
			},
		},
		{
			name:  "Genif with equality",
			input: "{%genif var == 2%}",
			expected: []parser.Node{
				&parser.StmtNode{
					StmtTagWithKw: parser.StmtTagWithKw{
						StmtTag: parser.StmtTag{
							PreWs: "",
						},
						Kw: parser.Kw{
							Kind: token.GENIF,
						},
					},
					Expr: &parser.BinaryExpr{
						X: &parser.Ident{
							Name: "var",
						},
						Op: parser.Kw{
							Kind: token.EQL,
						},
						Y: &parser.NumberLit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestSwitchStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "Switch with 1 case",
			input: `
{%switch name%}
{%case a%}
Text
{%end%}
`[1:],
			expected: []parser.Node{
				&parser.SwitchNode{
					SwitchTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
					},
					Cases: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "a",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"Text",
										"\n",
									},
								},
							},
						},
					},
					DefaultCase: nil,
				},
			},
		},
		{
			name: "Switch with several cases",
			input: `
{%switch name%}
{%case a%}
1
{%case b%}
2
{%case c%}
3
{%end%}
`[1:],
			expected: []parser.Node{
				&parser.SwitchNode{
					SwitchTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
					},
					Cases: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "a",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"1",
										"\n",
									},
								},
							},
						},
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "b",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"2",
										"\n",
									},
								},
							},
						},
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "c",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"3",
										"\n",
									},
								},
							},
						},
					},
					DefaultCase: nil,
				},
			},
		},
		{
			name: "Switch-default with 1 case",
			input: `
{%switch name%}
{%case a%}
Text
{%default%}
text2
{%end%}
`[1:],
			expected: []parser.Node{
				&parser.SwitchNode{
					SwitchTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
					},
					Cases: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "a",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"Text",
										"\n",
									},
								},
							},
						},
					},
					DefaultCase: &parser.Clause{
						Body: []parser.Node{
							&parser.TextNode{
								Val: []string{
									"text2",
									"\n",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Switch-default with several cases",
			input: `
{%switch name%}
{%case a%}
Text
{%case b%}
2
{%case c%}
3
{%default%}
text2
{%end%}
`[1:],
			expected: []parser.Node{
				&parser.SwitchNode{
					SwitchTag: parser.StmtTagWithExpr{
						Expr: &parser.Ident{
							Name: "name",
						},
					},
					Cases: []parser.ClauseWithExpr{
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "a",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"Text",
										"\n",
									},
								},
							},
						},
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "b",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"2",
										"\n",
									},
								},
							},
						},
						{
							Tag: parser.StmtTagWithExpr{
								Expr: &parser.Ident{
									Name: "c",
								},
							},
							Body: []parser.Node{
								&parser.TextNode{
									Val: []string{
										"3",
										"\n",
									},
								},
							},
						},
					},
					DefaultCase: &parser.Clause{
						Body: []parser.Node{
							&parser.TextNode{
								Val: []string{
									"text2",
									"\n",
								},
							},
						},
					},
				},
			},
		},
	}
	runTestCases(t, testCases)
}
