package analyzer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/types"
)

func TestGetTypeMap(t *testing.T) {
	testCases := []testCase{
		{
			name: "Plain text",
			input: parser.AST{
				&parser.Text{
					Value: "Hello world",
				},
			},
			expected: analyzer.TypeMap{},
		},
		{
			name: "Single var",
			input: parser.AST{
				&parser.Print{
					Expr: &parser.Ident{
						Name: "foo",
					},
				},
			},
			expected: analyzer.TypeMap{
				"foo": types.String,
			},
		},
		// 		{
		// 			name:  "Single var",
		// 			input: "{{ name }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 			},
		// 		},
		// 		{
		// 			name:  "String filter",
		// 			input: "{{ name -> upper }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var equal var",
		// 			input: "{{ name == surname }}",
		// 			expected: analyzer.TypeMap{
		// 				"name":    types.VarType("surname"),
		// 				"surname": types.Any,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var not equal var",
		// 			input: "{{ name != surname }}",
		// 			expected: analyzer.TypeMap{
		// 				"name":    types.VarType("surname"),
		// 				"surname": types.Any,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var is string",
		// 			input: "{{ name is 'asd' }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var and var",
		// 			input: "{{ name and b }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Boolean,
		// 				"b":    types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var equal number",
		// 			input: "{{ a == 1 }}",
		// 			expected: analyzer.TypeMap{
		// 				"a": types.Number,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var && var",
		// 			input: "{{ name && b }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Boolean,
		// 				"b":    types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "Vars infer boolean and then one of them used as string",
		// 			input: `
		// {{ name and b }}
		// {{ name }}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 				"b":    types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "Simple uniyfing of two vars to string",
		// 			input: `
		// {{ a == b }}
		// {{ a }}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"a": types.VarType("b"),
		// 				"b": types.String,
		// 			},
		// 		},
		// 		{
		// 			name: "Simple uniyfing of two vars to string",
		// 			input: `
		// {{ a == b }}
		// {% if a %}
		// {%end%}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"a": types.VarType("b"),
		// 				"b": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name:  "Vars and string",
		// 			input: "{{ name and 'asd' }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "Var infer boolean and then used as string",
		// 			input: `
		// {{ name and 'asd' }}
		// {{ name }}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 			},
		// 		},
		// 		{
		// 			name:  "Vars or number",
		// 			input: "{{ name or 1 }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Boolean,
		// 			},
		// 		},
		//
		// 		{
		// 			name:  "Parens changing precedence 1",
		// 			input: "{{(a or b) and c}}",
		// 			expected: analyzer.TypeMap{
		// 				"a": types.Boolean,
		// 				"b": types.Boolean,
		// 				"c": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var greater than var",
		// 			input: "{{ name > surname }}",
		// 			expected: analyzer.TypeMap{
		// 				"surname": types.Any,
		// 				"name":    types.VarType("surname"),
		// 			},
		// 		},
		// 		{
		// 			name:  "Var greater than number",
		// 			input: "{{ name > 1 }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Number,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var less than number",
		// 			input: "{{ name < 1 }}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.Number,
		// 			},
		// 		},
		// 		{
		// 			name:  "Var greater than string",
		// 			input: "{{ name > 'asd'}}",
		// 			expected: analyzer.TypeMap{
		// 				"name": types.String,
		// 			},
		// 		},
		// {
		// 	name: "Var + integer literal",
		// 	str:  "{{age+123}}",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Ident{
		// 					Name: "age",
		// 				},
		// 				Op: token.ADD,
		// 				Y: parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: analyzer.TypeMap{
		// 		"age": types.Number,
		// 	},
		// 	errExpected: false,
		// },
		// {
		// 	name: "Integer literal + var",
		// 	str:  "{{123+age}}",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.ADD,
		// 				Y: parser.Ident{
		// 					Name: "age",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: analyzer.TypeMap{
		// 		"age": types.Number,
		// 	},
		// 	errExpected: false,
		// },
		// {
		// 	name: "Var + var",
		// 	str:  "{{age+time}}",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Ident{
		// 					Name: "age",
		// 				},
		// 				Op: token.ADD,
		// 				Y: parser.Ident{
		// 					Name: "time",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: analyzer.TypeMap{
		// 		"age":  types.Any,
		// 		"time": types.Any,
		// 	},
		// 	errExpected: false,
		// },
		// {
		// 	name: "Var + string literal",
		// 	str:  "{{name+'ish'}}",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Ident{
		// 					Name: "name",
		// 				},
		// 				Op: token.ADD,
		// 				Y: parser.Lit{
		// 					Value: value.StringValue("ish"),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: analyzer.TypeMap{
		// 		"name": types.String,
		// 	},
		// 	errExpected: false,
		// },
		// 		{
		// 			name: "If statement on var",
		// 			input: `
		// {% if var %}
		// text
		// {% end %}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"var": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "If statement on var",
		// 			input: `
		// {% if var == 2 %}
		// text
		// {% end %}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"var": types.Number,
		// 			},
		// 		},
		// 		{
		// 			name: "Else-if statement on var",
		// 			input: `
		// {% if false %}
		// text
		// {% else if var %}
		// text2
		// {% end %}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"var": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "Nested If-else-if statements",
		// 			input: `
		// {% if var1 %}
		// text
		// {% else if var2 %}
		// {% if var3 %}
		// text123
		// {% else if var4 %}
		// 123
		// {% end %}
		// {% end %}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"var1": types.Boolean,
		// 				"var2": types.Boolean,
		// 				"var3": types.Boolean,
		// 				"var4": types.Boolean,
		// 			},
		// 		},
		// 		{
		// 			name: "Switch statement on var",
		// 			input: `
		// {% switch var %}
		// {% case 1 %}
		// {% case 2 %}
		// {% end %}
		// `[1:],
		// 			expected: analyzer.TypeMap{
		// 				"var": types.Number,
		// 			},
		// 		},
		// 		{
		// 			name: "Switch statement on var with different types",
		// 			input: `
		// {% switch var %}
		// {% case 1 %}
		// {% case "a" %}
		// {% case true %}
		// {% end %}
		// `[1:],
		// 			errExpected: analyzer.TypeErrors{
		// 				{
		// 					ExpectedType: types.String,
		// 					Name:         "var",
		// 				},
		// 				{
		// 					ExpectedType: types.Boolean,
		// 					Name:         "var",
		// 				},
		// 			},
		// 		},
	}
	runTestCases(t, testCases)
}
