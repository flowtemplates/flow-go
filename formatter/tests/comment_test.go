package renderer_test

import (
	"testing"
)

func TestCommentsUnchanged(t *testing.T) {
	testCases := []unchangedTestCase{
		{
			name: "Empty comment",
			input: `
{#  #}
`[1:],
		},
		{
			name: "Simple comment",
			input: `
{# comm #}
`[1:],
		},
		{
			name: "Comment inside text",
			input: `
Hello {# TODO #}.
text123123
`[1:],
		},
		{
			name: "Multiline comment",
			input: `
{# line 1
line 2
  
line 3 #}
`[1:],
		},
		{
			name: "Multiple comments",
			input: `
{# line 1 #}
{# line 2 #}
{# line 3 #}
`[1:],
		},
	}
	runUnchangedTestCases(t, testCases)
}

func TestComments(t *testing.T) {
	testCases := []testCase{
		{
			name: "Empty comment",
			input: `
{##}
`[1:],
			expected: `
{#  #}
`[1:],
		},
		{
			name: "Comment with spaces",
			input: `
{#  	#}
`[1:],
			expected: `
{#  #}
`[1:],
		},
		{
			name: "Simple comment",
			input: `
{# comm#}`[1:],
			expected: `
{# comm #}`[1:],
		},
		{
			name: "Multiline comment",
			input: `
{#line 1
line 2

line 3 #} `[1:],
			expected: `
{# line 1
line 2

line 3 #}`[1:],
		},
		{
			name: "Multiple comments",
			input: `
{#line 1 #}
{# 	 	line 2 #}
{# line 3#}
`[1:],
			expected: `
{# line 1 #}
{# line 2 #}
{# line 3 #}
`[1:],
		},
	}
	runTestCases(t, testCases)
}
