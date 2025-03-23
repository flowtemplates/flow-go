package lexer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/flowtemplates/flow-go/lexer"
	"github.com/flowtemplates/flow-go/token"
)

func TestPosition(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Empty if block",
			input: "{%if name%}\n{%end%}",
			expected: []token.Token{
				{Kind: token.LSTMT, Pos: token.Position{
					Line:   1,
					Column: 1,
					Offset: 0,
				}},
				{Kind: token.IF, Pos: token.Position{
					Line:   1,
					Column: 3,
					Offset: 2,
				}},
				{Kind: token.WS, Val: " ", Pos: token.Position{
					Line:   1,
					Column: 5,
					Offset: 4,
				}},
				{Kind: token.IDENT, Val: "name", Pos: token.Position{
					Line:   1,
					Column: 6,
					Offset: 5,
				}},
				{Kind: token.RSTMT, Pos: token.Position{
					Line:   1,
					Column: 10,
					Offset: 9,
				}},
				{Kind: token.TEXT, Val: "\n", Pos: token.Position{
					Line:   1,
					Column: 12,
					Offset: 11,
				}},
				{Kind: token.LSTMT, Pos: token.Position{
					Line:   2,
					Column: 1,
					Offset: 12,
				}},
				{Kind: token.END, Pos: token.Position{
					Line:   2,
					Column: 3,
					Offset: 14,
				}},
				{Kind: token.RSTMT, Pos: token.Position{
					Line:   2,
					Column: 6,
					Offset: 17,
				}},
			},
		},
		{
			name:  "Multiple expressions on many lines",
			input: "Hello {{name}}!\nFrom {{ flow }}templates\n\n{{3}}",
			expected: []token.Token{
				{Kind: token.TEXT, Val: "Hello ", Pos: token.Position{
					Line:   1,
					Column: 1,
					Offset: 0,
				}},
				{Kind: token.LEXPR, Pos: token.Position{
					Line:   1,
					Column: 7,
					Offset: 6,
				}},
				{Kind: token.IDENT, Val: "name", Pos: token.Position{
					Line:   1,
					Column: 9,
					Offset: 8,
				}},
				{Kind: token.REXPR, Pos: token.Position{
					Line:   1,
					Column: 13,
					Offset: 12,
				}},
				{Kind: token.TEXT, Val: "!", Pos: token.Position{
					Line:   1,
					Column: 15,
					Offset: 14,
				}},
				{Kind: token.LNBR, Val: "\n", Pos: token.Position{
					Line:   1,
					Column: 16,
					Offset: 15,
				}},
				{Kind: token.TEXT, Val: "From ", Pos: token.Position{
					Line:   2,
					Column: 1,
					Offset: 16,
				}},
				{Kind: token.LEXPR, Pos: token.Position{
					Line:   2,
					Column: 6,
					Offset: 21,
				}},
				{Kind: token.WS, Pos: token.Position{
					Line:   2,
					Column: 8,
					Offset: 23,
				}},
				{Kind: token.IDENT, Val: "flow", Pos: token.Position{
					Line:   2,
					Column: 9,
					Offset: 24,
				}},
				{Kind: token.WS, Pos: token.Position{
					Line:   2,
					Column: 13,
					Offset: 28,
				}},
				{Kind: token.REXPR, Pos: token.Position{
					Line:   2,
					Column: 14,
					Offset: 29,
				}},
				{Kind: token.TEXT, Val: "templates", Pos: token.Position{
					Line:   2,
					Column: 16,
					Offset: 31,
				}},
				{Kind: token.LNBR, Val: "\n", Pos: token.Position{
					Line:   2,
					Column: 25,
					Offset: 40,
				}},
				{Kind: token.LNBR, Val: "\n", Pos: token.Position{
					Line:   3,
					Column: 1,
					Offset: 41,
				}},
				{Kind: token.LEXPR, Pos: token.Position{
					Line:   4,
					Column: 1,
					Offset: 42,
				}},
				{Kind: token.INT, Val: "3", Pos: token.Position{
					Line:   4,
					Column: 3,
					Offset: 44,
				}},
				{Kind: token.REXPR, Pos: token.Position{
					Line:   4,
					Column: 4,
					Offset: 45,
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := lexer.TokensFromBytes([]byte(tc.input))

			eq := func(gotTokens []token.Token, expectedTokens []token.Token) error {
				l := len(gotTokens)
				if l != len(expectedTokens) {
					return errors.New("not matching length")
				}

				for i := range l {
					got, expected := gotTokens[i], expectedTokens[i]

					if expected.Pos != got.Pos {
						return fmt.Errorf("wrong pos: expected %+v, got %+v", expected.Pos, got.Pos)
					}
				}

				return nil
			}

			if err := eq(tokens, tc.expected); err != nil {
				t.Errorf("%s\nTest Case: %s\nExpected:\n%v\nGot:\n%v",
					err, tc.name, tc.expected, tokens)
			}
		})
	}
}
