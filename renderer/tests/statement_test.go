package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/renderer"
)

func TestIfStatements(t *testing.T) {
	testCases := []testCase{
		{
			name:     "Truthy if statement",
			input:    "{%if true%}\ntext\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "Falsy if statement",
			input:    "{%if false%}\ntext\n{%end%}",
			expected: "",
			scope:    renderer.Scope{},
		},
		{
			name:     "If with indentation",
			input:    "\t{%if true%}\n\ttext\n\t{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "If with end with not matching indentation level",
			input:    "{%if true%}\ntext\n\t{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "If with space indentation",
			input:    "  {%if true%}\n  text\n  {%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
		},
		{
			name:     "Nested ifs",
			input:    "{%if true%}\n\t{%if true%}\n\ttext\n\t{%end%}\n{%end%}",
			expected: "text\n",
			scope:    renderer.Scope{},
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
