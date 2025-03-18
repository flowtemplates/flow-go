package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/renderer"
)

type testCase struct {
	name        string
	input       string
	scope       renderer.Scope
	expected    string
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := renderer.RenderString(tc.input, tc.scope)
			if (err != nil) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.input, err)
				return
			}

			if got != tc.expected {
				t.Errorf("Input: %q\nMismatch.\nExpected:\n%q\nGot:\n%q", tc.input, tc.expected, got)
			}
		})
	}
}
