package analyzer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/types"
)

func TestGetTypeMap(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Plain text",
			input:    "Hello world",
			expected: analyzer.TypeMap{},
		},
		{
			name:  "Single var",
			input: "{{ name }}",
			expected: analyzer.TypeMap{
				"name": types.Any,
			},
		},
		{
			name:  "String filter",
			input: "{{ name -> upper }}",
			expected: analyzer.TypeMap{
				"name": types.String,
			},
		},
		{
			name:  "Var equal var",
			input: "{{ name == surname }}",
			expected: analyzer.TypeMap{
				"name":    types.Any,
				"surname": types.Any,
			},
			errExpected: false,
		},
		{
			name:  "Var greater than var",
			input: "{{ name > surname }}",
			expected: analyzer.TypeMap{
				"name":    types.Number,
				"surname": types.Number,
			},
			errExpected: false,
		},
		{
			name:  "Var greater than number",
			input: "{{ name > 1 }}",
			expected: analyzer.TypeMap{
				"name": types.Number,
			},
			errExpected: false,
		},
		{
			name:  "Var greater than string",
			input: "{{ name > 'asd'}}",
			expected: analyzer.TypeMap{
				"name": types.Number,
			},
			errExpected: false,
		},
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
		{
			name: "If statement on var",
			input: `
{% if var %}
text
{% end %}
`[1:],
			expected: analyzer.TypeMap{
				"var": types.Boolean,
			},
		},
		{
			name: "Else-if statement on var",
			input: `
{% if false %}
text
{% else if var %}
text2
{% end %}
`[1:],
			expected: analyzer.TypeMap{
				"var": types.Boolean,
			},
		},
		{
			name: "Nested If-else-if statements",
			input: `
{% if var1 %}
text
{% else if var2 %}
{% if var3 %}
text123
{% else if var4 %}
123
{% end %}
{% end %}
`[1:],
			expected: analyzer.TypeMap{
				"var1": types.Boolean,
				"var2": types.Boolean,
				"var3": types.Boolean,
				"var4": types.Boolean,
			},
		},
	}
	runTestCases(t, testCases)
}
