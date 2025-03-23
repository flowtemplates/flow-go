package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/renderer"
)

func TestIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "Truthy if statement",
			input: `
{%if true%}
text
{%end%}
`[1:],
			expected: `
text
`[1:],
			scope: renderer.Scope{},
		},
		{
			name: "Falsy if statement",
			input: `
{%if false%}
text
{%end%}
`[1:],
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name: "If with indentation",
			input: `
	{%if true%}
	text
	{%end%}
`[1:],
			expected: `
	text
`[1:],
			scope: renderer.Scope{},
		},
		{
			name: "If with end with not matching indentation level",
			input: `
{%if true%}
text
	{%end%}
`[1:],
			expected: `
text
`[1:],
			scope: renderer.Scope{},
		},
		{
			name: "If with space indentation",
			input: `
  {%if true%}
  text
  {%end%}
`[1:],
			expected: `
  text
`[1:],
			scope: renderer.Scope{},
		},
		{
			name: "If-else",
			input: `
{%if false%}
text
{%else%}
123
{%end%}
`[1:],
			expected: `
123
`[1:],
			scope: renderer.Scope{},
		},
		// 		{
		// 			name: "If-else-if",
		// 			input: `
		// {%if false%}
		// text
		// {%else if true%}
		// 123
		// {%end%}
		// `[1:],
		// 			expected: `
		// 123
		// `[1:],
		// 			scope: renderer.Scope{},
		// 		},
		{
			name: "Nested if-else statements",
			input: `
{% if true %}
function foo(n: number): number {
	{% if false%}
	return 2;
	{% else %}
	return 1;
	{% end %}
}
{% else %}
function bar() {
	console.log(123);
}
{% end %}
`[1:],
			expected: `
function foo(n: number): number {
	return 1;
}
`[1:],
			scope: renderer.Scope{},
		},
		{
			name:     "Not excl",
			input:    "{%if !false %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "Not",
			input:    "{%if not false %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
	}
	runTestCases(t, testCases)
}

func TestComparasions(t *testing.T) {
	testCases := []testCase{
		{
			name:     "String literals equal",
			input:    "{%if 'a'=='b'%}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "String literals is",
			input:    "{%if 'a' is 'b'%}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "String literals is not",
			input:    "{%if 'a' is not 'b'%}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "String literals not equal",
			input:    "{%if 'a' != 'b'%}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "Numbers greater",
			input:    "{%if 3 > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "3 < 2",
			input:    "{%if 3 < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "3 <= 2",
			input:    "{%if 3 <= 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "3 <= 3",
			input:    "{%if 3 <= 3 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "false < 2",
			input:    "{%if false < 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "'a' < 2",
			input:    "{%if 'a' < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		// TODO: convert string to NaN
		{
			name:     "'a' > 2",
			input:    "{%if 'a' > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "String greater than number",
			input:    "{%if '3' > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "String less than number",
			input:    "{%if '3' < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
	}
	runTestCases(t, testCases)
}
