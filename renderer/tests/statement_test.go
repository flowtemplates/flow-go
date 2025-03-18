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
	}
	runTestCases(t, testCases)
}
