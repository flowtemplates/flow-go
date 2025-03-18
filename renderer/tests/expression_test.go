package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/renderer"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:        "Plain text",
			input:       "Hello world",
			expected:    "Hello world",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Int literal",
			input:       "{{1}}",
			expected:    "1",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Float literal",
			input:       "{{1.1}}",
			expected:    "1.1",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Boolean literal",
			input:       "{{true}}",
			expected:    "",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "String literal",
			input:       `{{"word"}}`,
			expected:    "word",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		// {
		// 	name:     "Addition",
		// 	str:      "{{123+2}}",
		// 	expected: "125",
		// 	input: []parser.Node{
		// 		parser.ExprBlock{
		// 			Body: parser.BinaryExpr{
		// 				X: parser.Lit{
		// 					Value: value.NumberValue(123),
		// 				},
		// 				Op: token.ADD,
		// 				Y: parser.Lit{
		// 					Value: value.NumberValue(2),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	scope:       renderer.Scope{},
		// 	errExpected: false,
		// },
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
			input:    "{{name}}",
			expected: "useuse",
			scope: renderer.Scope{
				"name": "useuse",
			},
			errExpected: false,
		},
		{
			name:     "Multiple expressions",
			input:    "Hello {{name}}!\nFrom {{ flow }} templates",
			expected: "Hello world!\nFrom flow templates",
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
			name:        "Simple ternary with true",
			input:       "{{true?1:2}}",
			expected:    "1",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Simple ternary with false",
			input:       "{{false?1:2}}",
			expected:    "2",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Simple ternary with true and some text around",
			input:       "arr[{{false?1:2}}]",
			expected:    "arr[2]",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Simple ternary",
			input:    "{{flag?1:2}}",
			expected: "1",
			scope: renderer.Scope{
				"flag": true,
			},
			errExpected: false,
		},
		{
			name:     "Do-else ternary",
			input:    "{{flag do 1 else 2}}",
			expected: "1",
			scope: renderer.Scope{
				"flag": true,
			},
			errExpected: false,
		},
		{
			name:        "Ternary with truthy number condition",
			input:       "{{1?1:2}}",
			expected:    "1",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Ternary with falsy number condition",
			input:       "{{0?1:2}}",
			expected:    "2",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Ternary with truthy string condition",
			input:       `{{"a"?1:2}}`,
			expected:    "1",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:        "Ternary with falsy string condition",
			input:       `{{""?1:2}}`,
			expected:    "2",
			scope:       renderer.Scope{},
			errExpected: false,
		},
		{
			name:     "Ternary with 3 vars",
			input:    "{{flag?a:b}}",
			expected: "foo",
			scope: renderer.Scope{
				"flag": true,
				"a":    "foo",
				"b":    "bar",
			},
			errExpected: false,
		},
		{
			name:     "Ternary with truthy equal",
			input:    `{{flag==3?"foo":"bar"}}`,
			expected: "foo",
			scope: renderer.Scope{
				"flag": 3,
			},
			errExpected: false,
		},
		{
			name:     "Ternary with falsy equal",
			input:    `{{flag==4?"foo":"bar"}}`,
			expected: "bar",
			scope: renderer.Scope{
				"flag": 3,
			},
			errExpected: false,
		},
	}
	runTestCases(t, testCases)
}
