package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/parser"
	"github.com/flowtemplates/flow-go/renderer"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	name        string
	input       parser.AST
	scope       renderer.Input
	expected    string
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderer.RenderTemplate(tc.input, tc.scope)
			if (err != nil) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.input, err)

				return
			}

			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
