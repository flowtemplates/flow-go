package analyzer_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/parser"
)

type testCase struct {
	name        string
	input       parser.AST
	expected    analyzer.TypeMap
	errExpected analyzer.TypeErrors
}

func runTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tm, te := analyzer.TypeMapFromAST(tc.input)

			if diff := cmp.Diff(te, tc.errExpected); diff != "" {
				t.Errorf("errors mismatch (-got +want):\n%s", diff)
			}

			if diff := cmp.Diff(tm, tc.expected); diff != "" {
				t.Errorf("mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
