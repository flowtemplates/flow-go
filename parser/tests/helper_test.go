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
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parser.AstFromString(tc.input)

			if (err != nil) != tc.errExpected {
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
