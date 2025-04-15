package analyzer_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/flowtemplates/flow-go/analyzer"
)

type testCase struct {
	name        string
	input       string
	expected    analyzer.TypeMap
	errExpected analyzer.TypeErrors
}

func runTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := analyzer.New()

			err := a.TypeMapFromBytes([]byte(tc.input))

			switch {
			case tc.errExpected != nil:
				x, _ := json.MarshalIndent(tc.errExpected, "", "  ")
				y, _ := json.MarshalIndent(a.Errs, "", "  ")

				if (len(a.Errs) == 0) || !slices.Equal(x, y) {
					t.Errorf("Input: %q\nTypeMap mismatch.\nExpected:\n%q\nGot:\n%q", tc.input, x, y)
				}

			case err != nil:
				t.Errorf("Input: %q\nUnexpected parsing error: %v", tc.input, a.Errs)

			default:
				x, _ := json.MarshalIndent(tc.expected, "", "  ")
				y, _ := json.MarshalIndent(a.Tm, "", "  ")

				if !slices.Equal(x, y) {
					t.Errorf("Input: %q\nTypeMap mismatch.\nExpected:\n%s\nGot:\n%s", tc.input, x, y)
				}
			}
		})
	}
}
