package analyzer_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/flowtemplates/flow-go/analyzer"
	"github.com/flowtemplates/flow-go/parser"
)

type testCase struct {
	name        string
	str         string
	input       []parser.Node
	expected    analyzer.TypeMap
	errExpected bool
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make(analyzer.TypeMap)
			errs := analyzer.GetTypeMapFromAst(tc.input, got)
			if (len(errs) != 0) != tc.errExpected {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.str, errs)
				return
			}

			a, _ := json.MarshalIndent(tc.expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			if !slices.Equal(a, b) {
				t.Errorf("Input: %q\nTypeMap mismatch.\nExpected:\n%s\nGot:\n%s", tc.str, a, b)
			}
		})
	}
}
