package renderer_test

import (
	"testing"

	"github.com/flowtemplates/flow-go/formatter"
)

type testCase struct {
	name     string
	input    string
	expected string
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatter.FromBytes([]byte(tc.input))
			if err != nil {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.input, err)
				return
			}

			if string(got) != tc.expected {
				t.Errorf("Input: %q\nMismatch.\nExpected:\n%q\nGot:\n%q", tc.input, tc.expected, got)
			}
		})
	}
}

type unchangedTestCase struct {
	name  string
	input string
}

func runUnchangedTestCases(t *testing.T, testCases []unchangedTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatter.FromBytes([]byte(tc.input))
			if err != nil {
				t.Errorf("Input: %q\nUnexpected error: %v", tc.input, err)
				return
			}

			if string(got) != tc.input {
				t.Errorf("Input: %q\nMismatch.\nGot:\n%q", tc.input, got)
			}
		})
	}
}
