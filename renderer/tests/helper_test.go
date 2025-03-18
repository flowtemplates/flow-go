package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
)

type testCase struct {
	name        string
	str         string
	input       []parser.Node
	scope       renderer.Scope
	expected    string
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderer.RenderAst(tc.input, tc.scope)
			if (err != nil) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.str, err)
				return
			}

			if got != tc.expected {
				t.Errorf("Input: %q\nMismatch.\nExpected:\n%s\nGot:\n%s", tc.str, tc.expected, got)
			}
		})
	}
}
