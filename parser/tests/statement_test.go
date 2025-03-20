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
			name:  "Simple if statement",
			input: "{%if var%}\ntext\n{%end%}",
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
			name:  "If statement (with whitespaces)",
			input: "{% if var  %}\ntext\n{% end %}",
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
			name:  "If statement with indentation",
			input: "\t{%if var%}\n\ttext\n\t{%end%}",
			expected: []parser.Node{
				parser.IfStmt{
					BegTag: parser.StmtTagWithExpr{
						StmtTag: parser.StmtTag{
							PreWs: "\t",
							Kw:    token.IF,
						},
						Body: parser.Ident{
							Name: "var",
						},
					},
					Body: []parser.Node{
						parser.Text{
							Val: []string{
								"\t",
								"text",
								"\n",
							},
						},
					},
					Else:        nil,
					PreEndTagWs: "\t",
				},
			},
		},
		{
			name:  "Nested if blocks",
			input: "{%if var%}\n1\n{%if name%}\ntext\n{%end%}\n2\n{%end%}",
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

// func TestIfStatementsEdgeCases(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:        "If statement (with text before)",
// 			input:       "Text{%if var%}\ntext\n{%end%}",
// 			errExpected: true,
// 		},
// 	}
// 	runTestCases(t, testCases)
// }

func TestGenIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple genif",
			input: "{%genif var%}",
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
			name:  "Genif with equality",
			input: "{%genif var == 2%}",
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
	}
	runTestCases(t, testCases)
}

// func TestSwitchStatements(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name:  "Simple switch statement",
// 			input: "{%switch var%}\n{%case 1%}\ntext\n{%end%}",
// 			expected: []parser.Node{
// 				parser.SwitchStmt{
// 					BegTag: parser.StmtTagWithExpr{
// 						StmtTag: parser.StmtTag{
// 							Kw: token.SWITCH,
// 						},
// 						Body: parser.Ident{
// 							Name: "var",
// 						},
// 					},
// 					Cases: []parser.CaseClause{
// 						{
// 							CaseTag: parser.StmtTagWithExpr{
// 								StmtTag: parser.StmtTag{
// 									Kw: token.CASE,
// 								},
// 								Body: parser.Lit{
// 									Value: value.NumberValue(1),
// 								},
// 							},
// 							Body: []parser.Node{
// 								parser.Text{
// 									Val: []string{
// 										"text",
// 										"\n",
// 									},
// 								},
// 							},
// 						},
// 					},
// 					DefaultCase: []parser.Node{},
// 				},
// 				parser.IfStmt{
// 					BegTag: parser.StmtTagWithExpr{
// 						StmtTag: parser.StmtTag{
// 							Kw: token.IF,
// 						},
// 						Body: parser.Ident{
// 							Name: "var",
// 						},
// 					},
// 					Body: []parser.Node{
// 						parser.Text{
// 							Val: []string{
// 								"text",
// 								"\n",
// 							},
// 						},
// 					},
// 					Else: nil,
// 				},
// 			},
// 		},
// 	}
// 	runTestCases(t, testCases)
// }
