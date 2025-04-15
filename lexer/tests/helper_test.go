package lexer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/flowtemplates/flow-go/lexer"
	"github.com/flowtemplates/flow-go/token"
)

func equal(gotTokens []token.Token, expectedTokens []token.Token) error {
	l := len(gotTokens)
	if l != len(expectedTokens) {
		return errors.New("not matching length")
	}

	for i := range l {
		got, expected := gotTokens[i], expectedTokens[i]

		if got.Kind != expected.Kind {
			return fmt.Errorf("wrong type: expected %s, got %s", expected.Kind, got.Kind)
		}

		var expectedValue string
		if expected.IsValueable() {
			expectedValue = expected.Val
		} else {
			expectedValue = expected.Kind.String()
		}

		if got.Val != expectedValue {
			return fmt.Errorf("wrong value: expected %q, got %q", expectedValue, got.Val)
		}
	}

	return nil
}

type testCase struct {
	name     string
	input    string
	expected []token.Token
}

func runTestCases(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexer.TokensFromBytes([]byte(tc.input))

			if err := equal(tokens, tc.expected); err != nil {
				t.Errorf("%s\nInput: %s\nTest Case: %s\nExpected:\n%v\nGot:\n%v",
					err, tc.input, tc.name, tc.expected, tokens)
			}
		})
	}
}
