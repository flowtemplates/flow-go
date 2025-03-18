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
			str:  "{%if var%}\ntext\n{%end%}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							Kw: token.IF,
						},
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{
								"text",
								"\n",
							},
						},
					},
					Else: nil,
				},
			},
		},
		{
			name: "If statement (with whitespaces)",
			str:  "{% if var  %}text{% end %}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.WS, Val: " "},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: "  "},
				{Kind: token.RSTMT},
				{Kind: token.TEXT, Val: "text"},
				{Kind: token.LSTMT},
				{Kind: token.WS, Val: " "},
				{Kind: token.END},
				{Kind: token.WS, Val: " "},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							Kw: token.IF,
						},
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{"text"},
						},
					},
					Else: nil,
				},
			},
		},
		// {
		// 	name: "If statement (with indentations)",
		// 	str:  "\t{%if var%}\n\ttext\n\t{%end%}",
		// 	input: []token.Token{
		// 		{Typ: token.WS, Val: "\t"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.IF},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.IDENT, Val: "var"},
		// 		{Typ: token.RSTMT},
		// 		{Typ: token.LBR, Val: "\n"},
		// 		{Typ: token.WS, Val: "\t"},
		// 		{Typ: token.TEXT, Val: "text"},
		// 		{Typ: token.LBR, Val: "\n"},
		// 		{Typ: token.WS, Val: "\t"},
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.END},
		// 		{Typ: token.RSTMT},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.IfStmt{
		// 			PreLStmtBegWs:  "\t",
		// 			PostKwWs:       " ",
		// 			PostRStmtBegWs: "\n",
		// 			Condition: parser.Ident{
		// 				Name: "var",
		// 			},
		// 			Body: []parser.Node{
		// 				parser.Text{
		// 					Val: []string{
		// 						"\t",
		// 						"text",
		// 						"\n",
		// 					},
		// 				},
		// 			},
		// 			PostRStmtEndWs: "\t",
		// 			Else:           nil,
		// 		},
		// 	},
		// },
		{
			name: "If statement (with text before)",
			str:  "Text{%if var%}\ntext\n{%end%}",
			input: []token.Token{
				{Kind: token.TEXT, Val: "Text"},
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.Text{
					Val: []string{"Text"},
				},
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							Kw: token.IF,
						},
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{
								"text",
								"\n",
							},
						},
					},
					Else: nil,
				},
			},
		},
		{
			name: "Nested if blocks",
			str:  "{%if var%}\n1\n{%if name%}\ntext\n{%end%}\n2\n{%end%}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "1"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.IF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "name"},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "text"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.TEXT, Val: "2"},
				{Kind: token.LNBR, Val: "\n"},
				{Kind: token.LSTMT},
				{Kind: token.END},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							Kw: token.IF,
						},
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{"1", "\n"},
						},
						parser.IfStmt{
							BegTag: parser.StmtTagWithExpr{
								StmtTag: parser.StmtTag{
									Kw: token.IF,
								},
								Body: parser.Ident{
									Name: "name",
								},
							},
							Body: []parser.Node{
								parser.Text{
									Val: []string{"text", "\n"},
								},
							},
						},
						parser.Text{
							Val: []string{"2", "\n"},
						},
					},
					Else: nil,
				},
			},
		},
	}
	runTestCases(t, testCases)
}

func TestGenIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "Simple genif",
			str:  "{%genif var%}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.GENIF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.StmtTagWithExpr{
					StmtTag: parser.StmtTag{
						Kw: token.GENIF,
					},
					Body: parser.Ident{
						Name: "var",
					},
				},
			},
		},
		{
			name: "Genif with equality",
			str:  "{%genif var == 2%}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.GENIF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: " "},
				{Kind: token.EQL},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "2"},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.StmtTagWithExpr{
					StmtTag: parser.StmtTag{
						Kw: token.GENIF,
					},
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "var",
						},
						Op: token.EQL,
						Y: parser.Lit{
							Value: value.NumberValue(2),
						},
					},
				},
			},
		},
		{
			name: "Genif with comparison",
			str:  "{%genif var > 3%}",
			input: []token.Token{
				{Kind: token.LSTMT},
				{Kind: token.GENIF},
				{Kind: token.WS, Val: " "},
				{Kind: token.IDENT, Val: "var"},
				{Kind: token.WS, Val: " "},
				{Kind: token.GTR},
				{Kind: token.WS, Val: " "},
				{Kind: token.INT, Val: "3"},
				{Kind: token.RSTMT},
			},
			expected: []parser.Node{
				parser.StmtTagWithExpr{
					StmtTag: parser.StmtTag{
						Kw: token.GENIF,
					},
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "var",
						},
						Op: token.GTR,
						Y: parser.Lit{
							Value: value.NumberValue(3),
						},
					},
				},
			},
		},
		// TODO
		// {
		// 	name: "Genif with complex expression",
		// 	str:  "{%genif (var+1)/2+name%}",
		// 	input: []token.Token{
		// 		{Typ: token.LSTMT},
		// 		{Typ: token.GENIF},
		// 		{Typ: token.WS, Val: " "},
		// 		{Typ: token.LPAREN},
		// 		{Typ: token.IDENT, Val: "var"},
		// 		{Typ: token.ADD},
		// 		{Typ: token.INT, Val: "1"},
		// 		{Typ: token.RPAREN},
		// 		{Typ: token.DIV},
		// 		{Typ: token.INT, Val: "2"},
		// 		{Typ: token.ADD},
		// 		{Typ: token.IDENT, Val: "name"},
		// 		{Typ: token.RSTMT},
		// 	},
		// 	expected: []parser.Node{
		// 		parser.GenIfStmt{
		// 			PostKwWs: " ",
		// 			Condition: parser.BinaryExpr{
		// 				X: parser.BinaryExpr{
		// 					X: parser.BinaryExpr{
		// 						X: parser.Ident{
		// 							Name: "var",
		// 						},
		// 						PostOpWs: "",
		// 						Op:       token.ADD,
		// 						Y: parser.Lit{
		// 							Typ: token.INT,
		// 							Val: "1",
		// 						},
		// 					},
		// 					Op: token.DIV,
		// 					Y: parser.Lit{
		// 						Typ: token.INT,
		// 						Val: "2",
		// 					},
		// 				},
		// 				PostOpWs: "",
		// 				Op:       token.ADD,
		// 				Y: parser.Ident{
		// 					Name: "name",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
	runTestCases(t, testCases)
}
