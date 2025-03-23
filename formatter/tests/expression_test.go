package renderer_test

import (
	"testing"
)

func TestTextUnchanged(t *testing.T) {
	testCases := []unchangedTestCase{
		{
			name: "Plain text",
			input: `
Hello world
`[1:],
		},
		{
			name: "Plain text 2",
			input: `
Hello world
`[1:],
		},
		{
			name: "Plain text 3",
			input: `
Hello world
123
asdasd
`[1:],
		},
	}
	runUnchangedTestCases(t, testCases)
}

func TestExpressionsUnchanged(t *testing.T) {
	testCases := []unchangedTestCase{
		{
			name: "Expression with int literal",
			input: `
{{ 1 }}
`[1:],
		},
		{
			name: "Expression with negative int lit",
			input: `
{{ -1 }}
`[1:],
		},
		{
			name: "Expression with float lit",
			input: `
{{ 1.1 }}
`[1:],
		},
		{
			name: "Expression with negative float lit",
			input: `
{{ -1.1 }}
`[1:],
		},
		{
			name: "Expression with single quote string lit",
			input: `
{{ '1' }}
`[1:],
		},
		{
			name: "Expression with double quote string lit",
			input: `
{{ "1" }}
`[1:],
		},
		{
			name: "True",
			input: `
{{ true }}
`[1:],
		},
		{
			name: "Ident",
			input: `
{{ var }}
`[1:],
		},
		{
			name: "Equal",
			input: `
{{ 1 == 2 }}
`[1:],
		},
		{
			name: "Not equal",
			input: `
{{ 1 != 2 }}
`[1:],
		},
		{
			name: "Is",
			input: `
{{ 1 is 2 }}
`[1:],
		},
		{
			name: "Not var",
			input: `
{{ not var }}
`[1:],
		},
		{
			name: "Excl var",
			input: `
{{ !var }}
`[1:],
		},
		{
			name: "a is not b",
			input: `
{{ a is not b }}
`[1:],
		},
		{
			name: "Greater",
			input: `
{{ 1 > 2 }}
`[1:],
		},
		{
			name: "Less",
			input: `
{{ 1 < 2 }}
`[1:],
		},
		{
			name: "Greater or equal",
			input: `
{{ 1 >= 2 }}
`[1:],
		},
		{
			name: "Less or equal",
			input: `
{{ 1 <= 2 }}
`[1:],
		},
		{
			name: "And",
			input: `
{{ 1 and 2 }}
`[1:],
		},
		{
			name: "Or",
			input: `
{{ 1 or 2 }}
`[1:],
		},
		{
			name: "&&",
			input: `
{{ 1 && 2 }}
`[1:],
		},
		{
			name: "||",
			input: `
{{ 1 || 2 }}
`[1:],
		},
		{
			name: "Redundant parens",
			input: `
{{ (a) }}
`[1:],
		},
		{
			name: "Logical expression with parens",
			input: `
{{ (var || 1) && false }}
`[1:],
		},
		{
			name: "Multiple expressions",
			input: `
{{ 2 }}+{{ 1 }}
{{ -123 }}
`[1:],
		},
		{
			name: "Multiple expressions with text in between",
			input: `
{{ 2 }}
text
{{ 1 }}
text
{{ -123 }}
`[1:],
		},
		{
			name: "Simple ternary",
			input: `
{{ flag ? 1 : 2 }}
`[1:],
		},
		{
			name: "Simple do-else ternary",
			input: `
{{ flag do 1 else 2 }}
`[1:],
		},
		{
			name: "Simple ternary with equal",
			input: `
{{ asd == 2 ? 1 : 2 }}
`[1:],
		},
		{
			name: "Simple do-else ternary with equal",
			input: `
{{ asd == 2 do 1 else 2 }}
`[1:],
		},
		{
			name: "Nested ternaries",
			input: `
{{ flag ? bar ? 1 : 3 : 2 }}
`[1:],
		},
		{
			name: "Nested do-else ternaries",
			input: `
{{ flag do bar ? 1 : 3 else 2 }}
`[1:],
		},
	}
	runUnchangedTestCases(t, testCases)
}

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			name: "Expression with int literal",
			input: `
{{1}}
`[1:],
			expected: `
{{ 1 }}
`[1:],
		},
		{
			name: "Expression with negative int lit",
			input: `
{{ -1}}
`[1:],
			expected: `
{{ -1 }}
`[1:],
		},
		{
			name: "Expression with float lit",
			input: `
{{1.1 }}
`[1:],
			expected: `
{{ 1.1 }}
`[1:],
		},
		{
			name: "Expression with negative float lit",
			input: `
{{ -1.1}}
`[1:],
			expected: `
{{ -1.1 }}
`[1:],
		},
		{
			name: "Equal",
			input: `
{{1==2}}
`[1:],
			expected: `
{{ 1 == 2 }}
`[1:],
		},
		{
			name: "Is",
			input: `
{{1  is 	2}}
`[1:],
			expected: `
{{ 1 is 2 }}
`[1:],
		},
		{
			name: "Excl var",
			input: `
{{ ! var}}
`[1:],
			expected: `
{{ !var }}
`[1:],
		},
		{
			name: "Redundant parens",
			input: `
{{(a )}}
`[1:],
			expected: `
{{ (a) }}
`[1:],
		},
		{
			name: "Multiple expressions with text in between",
			input: `
{{2}}
text
{{1}}
text
{{-123 }}
`[1:],
			expected: `
{{ 2 }}
text
{{ 1 }}
text
{{ -123 }}
`[1:],
		},
		{
			name: "Simple ternary",
			input: `
{{flag?1:2}}
`[1:],
			expected: `
{{ flag ? 1 : 2 }}
`[1:],
		},
		{
			name: "Simple do-else ternary",
			input: `
{{flag do 1else 2}}
`[1:], // TODO: maybe do not allow this '1else' thing
			expected: `
{{ flag do 1 else 2 }}
`[1:],
		},
		{
			name: "Simple ternary with equal",
			input: `
{{asd 	==2 ?1 :  2  }}
`[1:],
			expected: `
{{ asd == 2 ? 1 : 2 }}
`[1:],
		},
		{
			name: "Nested ternaries",
			input: `
{{flag ?	 bar? 1 :3 :2 }}
`[1:],
			expected: `
{{ flag ? bar ? 1 : 3 : 2 }}
`[1:],
		},
	}
	runTestCases(t, testCases)
}
