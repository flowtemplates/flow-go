package analyzer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/token"
	"github.com/flowtemplates/flow-go/types"
	"github.com/flowtemplates/flow-go/value"
)

func TestGetTypeMap(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			str:  "Hello world",
			input: []parser.Node{
				parser.Text{
					Val: []string{"Hello world"},
				},
			},
			expected:    analyzer.TypeMap{},
			errExpected: false,
		},
		{
			name: "Single var",
			str:  "{{name}}",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.Ident{Name: "name"},
				},
			},
			expected: analyzer.TypeMap{
				"name": types.Any,
			},
			errExpected: false,
		},
		{
			name: "Var + integer literal",
			str:  "{{age+123}}",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "age",
						},
						Op: token.ADD,
						Y: parser.Lit{
							Value: value.NumberValue(123),
						},
					},
				},
			},
			expected: analyzer.TypeMap{
				"age": types.Number,
			},
			errExpected: false,
		},
		{
			name: "Integer literal + var",
			str:  "{{123+age}}",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Lit{
							Value: value.NumberValue(123),
						},
						Op: token.ADD,
						Y: parser.Ident{
							Name: "age",
						},
					},
				},
			},
			expected: analyzer.TypeMap{
				"age": types.Number,
			},
			errExpected: false,
		},
		{
			name: "Var + var",
			str:  "{{age+time}}",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "age",
						},
						Op: token.ADD,
						Y: parser.Ident{
							Name: "time",
						},
					},
				},
			},
			expected: analyzer.TypeMap{
				"age":  types.Any,
				"time": types.Any,
			},
			errExpected: false,
		},
		{
			name: "Var + string literal",
			str:  "{{name+'ish'}}",
			input: []parser.Node{
				parser.ExprBlock{
					Body: parser.BinaryExpr{
						X: parser.Ident{
							Name: "name",
						},
						Op: token.ADD,
						Y: parser.Lit{
							Value: value.StringValue("ish"),
						},
					},
				},
			},
			expected: analyzer.TypeMap{
				"name": types.String,
			},
			errExpected: false,
		},
		{
			name: "If statement",
			str:  "{%if var%}\ntext\n{%end%}",
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
			expected: analyzer.TypeMap{
				"var": types.Boolean,
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
