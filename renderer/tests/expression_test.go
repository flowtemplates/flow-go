package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/renderer"
)

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Plain text",
			input:    "Hello world",
			expected: "Hello world",
			scope:    renderer.Scope{},
		},
		{
			name:     "Int literal",
			input:    "{{1}}",
			expected: "1",
			scope:    renderer.Scope{},
		},
		{
			name:     "Float literal",
			input:    "{{1.1}}",
			expected: "1.1",
			scope:    renderer.Scope{},
		},
		{
			name:     "Boolean literal",
			input:    "{{true}}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "String literal in double quotes",
			input:    `{{"word"}}`,
			expected: "word",
			scope:    renderer.Scope{},
		},
		{
			name:     "String literal in single quotes",
			input:    `{{'word'}}`,
			expected: "word",
			scope:    renderer.Scope{},
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
		// },
		{
			name:     "Expression with string var",
			input:    "{{name}}",
			expected: "useuse",
			scope: renderer.Scope{
				"name": "useuse",
			},
		},
		{
			name:     "Expression with number var",
			input:    "{{age}}",
			expected: "1",
			scope: renderer.Scope{
				"age": 1,
			},
		},
		{
			name:     "Expression with boolean var",
			input:    "{{flag}}",
			expected: "",
			scope: renderer.Scope{
				"flag": false,
			},
		},
		{
			name: "Multiple expressions",
			input: `
Hello {{name}}!
From {{ flow }} templates
`[1:],
			expected: `
Hello world!
From flow templates
`[1:],
			scope: renderer.Scope{
				"name": "world",
				"flow": "flow",
			},
		},
	}
	runTestCases(t, testCases)
}

func TestOperators(t *testing.T) {
	testCases := []testCase{
		{
			name:     "String literals and",
			input:    "{{'a' && 'b'}}",
			expected: "b",
			scope:    renderer.Scope{},
		},
		{
			name:     "Number literals and",
			input:    "{{1 && 2}}",
			expected: "2",
			scope:    renderer.Scope{},
		},
		{
			name:     "Number literals falsy and",
			input:    "{{0 && 0}}",
			expected: "0",
			scope:    renderer.Scope{},
		},
		{
			name:     "Equality with empty string",
			input:    "{{0 == ''}}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "Number literals or",
			input:    "{{1 || 2}}",
			expected: "1",
			scope:    renderer.Scope{},
		},
		{
			name:     "0 || 'a'",
			input:    "{{0 || 'a'}}",
			expected: "a",
			scope:    renderer.Scope{},
		},
		{
			name:     "Multiple || with strings",
			input:    "{{'a' || 'b' || 'c'}}",
			expected: "a",
			scope:    renderer.Scope{},
		},
		{
			name:     "Multiple && with strings",
			input:    "{{'a' && 'b' && 'c'}}",
			expected: "c",
			scope:    renderer.Scope{},
		},
	}
	runTestCases(t, testCases)
}

func TestTernaries(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Simple ternary with true",
			input:    "{{true?1:2}}",
			expected: "1",
			scope:    renderer.Scope{},
		},
		{
			name:     "Simple ternary with false",
			input:    "{{false?1:2}}",
			expected: "2",
			scope:    renderer.Scope{},
		},
		{
			name:     "Simple ternary with true and some text around",
			input:    "arr[{{false?1:2}}]",
			expected: "arr[2]",
			scope:    renderer.Scope{},
		},
		{
			name:     "Simple ternary",
			input:    "{{flag?1:2}}",
			expected: "1",
			scope: renderer.Scope{
				"flag": true,
			},
		},
		{
			name:     "Do-else ternary",
			input:    "{{flag do 1 else 2}}",
			expected: "1",
			scope: renderer.Scope{
				"flag": true,
			},
		},
		{
			name:     "Ternary with truthy number condition",
			input:    "{{1?1:2}}",
			expected: "1",
			scope:    renderer.Scope{},
		},
		{
			name:     "Ternary with falsy number condition",
			input:    "{{0?1:2}}",
			expected: "2",
			scope:    renderer.Scope{},
		},
		{
			name:     "Ternary with truthy string condition",
			input:    `{{"a"?1:2}}`,
			expected: "1",
			scope:    renderer.Scope{},
		},
		{
			name:     "Ternary with falsy string condition",
			input:    `{{""?1:2}}`,
			expected: "2",
			scope:    renderer.Scope{},
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
		},
		{
			name:     "Ternary with truthy equal",
			input:    `{{flag==3?"foo":"bar"}}`,
			expected: "foo",
			scope: renderer.Scope{
				"flag": 3,
			},
		},
		{
			name:     "Ternary with falsy equal",
			input:    `{{flag==4?"foo":"bar"}}`,
			expected: "bar",
			scope: renderer.Scope{
				"flag": 3,
			},
		},
	}
	runTestCases(t, testCases)
}

func TestFilters(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Empty string upper",
			input:    "{{ '' -> upper }}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "String with emoji upper",
			input:    "{{ '💀' -> upper }}",
			expected: "💀",
			scope:    renderer.Scope{},
		},
		{
			name:     "Upper",
			input:    "{{ s -> upper }}",
			expected: "HELLO WORLD",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Upper to string lit",
			input:    "{{ 'Hello world' -> upper }}",
			expected: "HELLO WORLD",
			scope:    renderer.Scope{},
		},
		{
			name:     "Upper to number lit",
			input:    "{{ 123 -> upper }}",
			expected: "123",
			scope:    renderer.Scope{},
		},
		{
			name:     "Upper to true",
			input:    "{{ true -> upper }}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "Lower",
			input:    "{{ s -> lower }}",
			expected: "hello world",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Camel case",
			input:    "{{ s -> camel }}",
			expected: "helloWorld",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Pascal case",
			input:    "{{ s -> pascal }}",
			expected: "HelloWorld",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Kebab case",
			input:    "{{ s -> kebab }}",
			expected: "hello-world",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Snake case",
			input:    "{{ s -> snake }}",
			expected: "hello_world",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Title case",
			input:    "{{ s -> title }}",
			expected: "Hello World",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Capitalize",
			input:    "{{ s -> capitalize}}",
			expected: "Hello world",
			scope: renderer.Scope{
				"s": "hello world",
			},
		},
		{
			name:     "Trim",
			input:    "{{ s -> trim }}",
			expected: "Hello world",
			scope: renderer.Scope{
				"s": "  Hello world 	",
			},
		},
		{
			name:     "String length",
			input:    "{{ s -> length }}",
			expected: "11",
			scope: renderer.Scope{
				"s": "Hello world",
			},
		},
		{
			name:     "Number length",
			input:    "{{ 123 -> length }}",
			expected: "3",
			scope:    renderer.Scope{},
		},
		{
			name:     "Boolean length",
			input:    "{{ true -> length }}",
			expected: "0",
			scope:    renderer.Scope{},
		},
		{
			name:     "Var name 'length'",
			input:    "{{ length -> length }}",
			expected: "3",
			scope: renderer.Scope{
				"length": "huh",
			},
		},
	}
	runTestCases(t, testCases)
}
