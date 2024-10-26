package tokenizer

import (
	"reflect"
	"testing"

	"github.com/LitoleNINJA/json-parser/cmd/tokenizer"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input      []byte
		expected   []tokenizer.Token
		shouldFail bool
	}{
		{
			input: []byte(`{"key": "value"}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "key"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.String, Value: "value"},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"outer": {"inner": "value"}}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "outer"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "inner"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.String, Value: "value"},
				{Type: tokenizer.RightBrace, Value: "}"},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`["apple", "banana", "cherry"]`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBracket, Value: "["},
				{Type: tokenizer.String, Value: "apple"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "banana"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "cherry"},
				{Type: tokenizer.RightBracket, Value: "]"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`[true, null, 42, "text"]`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBracket, Value: "["},
				{Type: tokenizer.Boolean, Value: "true"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.Boolean, Value: "null"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.Number, Value: "42"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "text"},
				{Type: tokenizer.RightBracket, Value: "]"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"isValid": true, "isNull": null, "isFalse": false}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "isValid"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Boolean, Value: "true"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "isNull"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Boolean, Value: "null"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "isFalse"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Boolean, Value: "false"},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"zero": 0, "neg": -10, "float": 10.5}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "zero"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Number, Value: "0"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "neg"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Number, Value: "-10"},
				{Type: tokenizer.Comma, Value: ","},
				{Type: tokenizer.String, Value: "float"},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.Number, Value: "10.5"},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"   key   ": "   value   "}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.String, Value: "   key   "},
				{Type: tokenizer.Colon, Value: ":"},
				{Type: tokenizer.String, Value: "   value   "},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input:      []byte(`{"key": "value"`), // Invalid case
			expected:   nil,
			shouldFail: true,
		},
		{
			input: []byte(`{}`),
			expected: []tokenizer.Token{
				{Type: tokenizer.LeftBrace, Value: "{"},
				{Type: tokenizer.RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input:      []byte(`{"key": "value"} // This is a comment`),
			expected:   nil,
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			tokens, err := tokenizer.TokenizeJSON(tc.input)

			if tc.shouldFail {
				if err == nil {
					t.Fatalf("Expected an error for input %s, but got none.", tc.input)
				}
				return // Skip further checks if we expect failure
			}

			if err != nil {
				t.Fatalf("Tokenize(%s) returned an error: %v", tc.input, err)
			}

			if !reflect.DeepEqual(tokens, tc.expected) {
				t.Errorf("Tokenize(%s) = %v, want %v", tc.input, tokens, tc.expected)
			}
		})
	}
}
