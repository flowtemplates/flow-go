package renderer_test

import (
	"testing"
)

func TestStatementsUnchanged(t *testing.T) {
	testCases := []unchangedTestCase{
		{
			name: "Simple if statement",
			input: `
{% if var %}
text
{% end %}
`[1:],
		},
		{
			name: "If-else statement",
			input: `
{% if var %}
text
{% else %}
2
{% end %}
`[1:],
		},
		{
			name: "If-elseif statement",
			input: `
{% if bar %}
1
{% else if flag %}
2
{% end %}
`[1:],
		},
		{
			name: "If-elseif-else statement",
			input: `
{% if bar %}
1
{% else if flag %}
2
{% else %}
3
{% end %}
`[1:],
		},
		{
			name: "Nested if statements",
			input: `
{% if var %}
1
{% if name %}
text
{% end %}
2
{% end %}
`[1:],
		},
		{
			name: "If with indentation",
			input: `
	{% if var %}
	2
	{% end %}
`[1:],
		},
		{
			name: "If-else statement with indentation",
			input: `
	{% if var %}
	text
	{% else %}
	2
	{% end %}
`[1:],
		},
		// TODO: indent all tags according to if tag (?)
		{
			name: "If-else statement with various indentations",
			input: `
	{% if var %}
	text
		{% else %}
	2
  {% end %}
`[1:],
		},
		{
			name: "Switch",
			input: `
{% switch name %}
{% case a %}
Text
{% case b %}
123
{% end %}
`[1:],
		},
		{
			name: "Switch-default",
			input: `
{% switch name %}
{% case a %}
Text
{% case b %}
123
{% default %}
456
{% end %}
`[1:],
		},
	}
	runUnchangedTestCases(t, testCases)
}

func TestStatements(t *testing.T) {
	testCases := []testCase{
		{
			name: "If statement with whitespaces after tag",
			input: `
{% if var%}   
text
{% end%}	
`[1:],
			expected: `
{% if var %}
text
{% end %}
`[1:],
		},
		{
			name: "Switch",
			input: `
{%switch name%}
{%case a   %}
Text
{% end%} 
`[1:],
			expected: `
{% switch name %}
{% case a %}
Text
{% end %}
`[1:],
		},
		{
			name: "Switch-default",
			input: `
{% switch name %}
{%case a %}
Text
{%   default %}
456
{% end %}
`[1:],
			expected: `
{% switch name %}
{% case a %}
Text
{% default %}
456
{% end %}
`[1:],
		},
	}
	runTestCases(t, testCases)
}
