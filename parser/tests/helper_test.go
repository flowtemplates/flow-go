package parser_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/flowtemplates/flow-go/parser"
)

type testCase struct {
	name        string
	input       string
	expected    parser.Expr
	errExpected error
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.AstFromString(tc.input)

			if tc.errExpected != nil {
				if tc.errExpected.Error() != err.Error() {
					t.Errorf("Input: %q\nUnexpected error: %v, got: %v", tc.input, err, tc.errExpected)
					return
				}
				return
			} else if err != nil {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.input, err)
				return
			}

			a, _ := json.MarshalIndent(tc.expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			if !slices.Equal(a, b) {
				t.Errorf("Input: %q\nAST mismatch.\nExpected:\n%s\nGot:\n%s", tc.input, a, b)
			}
		})
	}
}
