package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		input      []byte
		expected   []Token
		shouldFail bool
	}{
		{
			input: []byte(`{"key": "value"}`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "key"},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "value"},
				{Type: RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"outer": {"inner": "value"}}`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "outer"},
				{Type: Colon, Value: ":"},
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "inner"},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "value"},
				{Type: RightBrace, Value: "}"},
				{Type: RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`["apple", "banana", "cherry"]`),
			expected: []Token{
				{Type: LeftBracket, Value: "["},
				{Type: String, Value: "apple"},
				{Type: Comma, Value: ","},
				{Type: String, Value: "banana"},
				{Type: Comma, Value: ","},
				{Type: String, Value: "cherry"},
				{Type: RightBracket, Value: "]"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`[true, null, 42, "text"]`),
			expected: []Token{
				{Type: LeftBracket, Value: "["},
				{Type: Boolean, Value: "true"},
				{Type: Comma, Value: ","},
				{Type: Null, Value: "null"},
				{Type: Comma, Value: ","},
				{Type: Number, Value: "42"},
				{Type: Comma, Value: ","},
				{Type: String, Value: "text"},
				{Type: RightBracket, Value: "]"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"isValid": true, "isNull": null, "isFalse": false}`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "isValid"},
				{Type: Colon, Value: ":"},
				{Type: Boolean, Value: "true"},
				{Type: Comma, Value: ","},
				{Type: String, Value: "isNull"},
				{Type: Colon, Value: ":"},
				{Type: Null, Value: "null"},
				{Type: Comma, Value: ","},
				{Type: String, Value: "isFalse"},
				{Type: Colon, Value: ":"},
				{Type: Boolean, Value: "false"},
				{Type: RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"   key   ": "   value   "}`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "   key   "},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "   value   "},
				{Type: RightBrace, Value: "}"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"key": "value"`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: String, Value: "key"},
				{Type: Colon, Value: ":"},
				{Type: String, Value: "value"},
			},
			shouldFail: false,
		},
		{
			input: []byte(`{}`),
			expected: []Token{
				{Type: LeftBrace, Value: "{"},
				{Type: RightBrace, Value: "}"},
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
			tokens, err := Tokenize(tc.input)

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
