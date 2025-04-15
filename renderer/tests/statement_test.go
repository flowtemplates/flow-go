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
			scope: renderer.Input{},
		},
		{
			name: "Falsy if statement",
			input: `
{%if false%}
text
{%end%}
`[1:],
			expected: "",
			scope:    renderer.Input{},
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
			scope: renderer.Input{},
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
			scope: renderer.Input{},
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
			scope: renderer.Input{},
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
			scope: renderer.Input{},
		},
		{
			name: "If-else-if",
			input: `
{%if false%}
text
{%else if true%}
123
{%end%}
`[1:],
			expected: `
123
`[1:],
			scope: renderer.Input{},
		},
		{
			name: "Nested If-else-if statements",
			input: `
{% if false %}
text
{% else if true %}
{% if false %}
text123
{% else if true %}
123
{% end %}
{% end %}
`[1:],
			expected: `
123
`[1:],
			scope: renderer.Input{},
		},
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
			scope: renderer.Input{},
		},
		{
			name:     "Not excl",
			input:    "{%if !false %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "Not",
			input:    "{%if not false %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
	}
	runTestCases(t, testCases)
}

func TestSwitchStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "Switch with 1 case",
			input: `
{% switch a %}
{% case 1 %}
Text
{% end %}
`[1:],
			expected: `
Text
`[1:],
			scope: renderer.Input{
				"a": 1,
			},
		},
		{
			name: "False switch with 1 case",
			input: `
{% switch a %}
{% case 2 %}
Text
{% end %}
`[1:],
			expected: `
`[1:],
			scope: renderer.Input{
				"a": 1,
			},
		},
		{
			name: "Switch with several cases",
			input: `
{% switch a %}
{% case 1 %}
111
{% case 2 %}
22
{% case 4 %}
33
{% end %}
`[1:],
			expected: `
22
`[1:],
			scope: renderer.Input{
				"a": 2,
			},
		},
		{
			name: "Switch-default with 1 case",
			input: `
{% switch a %}
{% case 1 %}
Text
{% default %}
Text2
{% end %}
`[1:],
			expected: `
Text2
`[1:],
			scope: renderer.Input{
				"a": 2,
			},
		},
		{
			name: "Switch on string with number",
			input: `
{% switch 'foo' %}
{% case 1 %}
Text
{% end %}
`[1:],
			expected: `
`[1:],
			scope: renderer.Input{},
		},
		{
			name: "Switch on string convertable to number with number",
			input: `
{% switch '2' %}
{% case 1 %}
Text
{% case 2 %}
text2
{% end %}
`[1:],
			expected: `
text2
`[1:],
			scope: renderer.Input{},
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
			scope:    renderer.Input{},
		},
		{
			name:     "String literals is",
			input:    "{%if 'a' is 'b'%}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Input{},
		},
		{
			name:     "String literals is not",
			input:    "{%if 'a' is not 'b'%}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "String literals not equal",
			input:    "{%if 'a' != 'b'%}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "Numbers greater",
			input:    "{%if 3 > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "3 < 2",
			input:    "{%if 3 < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Input{},
		},
		{
			name:     "3 <= 2",
			input:    "{%if 3 <= 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Input{},
		},
		{
			name:     "3 <= 3",
			input:    "{%if 3 <= 3 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "false < 2",
			input:    "{%if false < 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "'a' < 2",
			input:    "{%if 'a' < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Input{},
		},
		// TODO: convert string to NaN
		{
			name:     "'a' > 2",
			input:    "{%if 'a' > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "String greater than number",
			input:    "{%if '3' > 2 %}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Input{},
		},
		{
			name:     "String less than number",
			input:    "{%if '3' < 2 %}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Input{},
		},
	}
	runTestCases(t, testCases)
}
