package parser

import (
	"reflect"
	"testing"

	"github.com/LitoleNINJA/json-parser/cmd/tokenizer"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		input      []byte
		expected   interface{}
		shouldFail bool
	}{
		{
			input: []byte(`{"name": "John", "age": 30, "isStudent": true}`),
			expected: JsonObject{
				"name":      "John",
				"age":       int64(30),
				"isStudent": true,
			},
			shouldFail: false,
		},
		{
			input: []byte(`["apple", "banana", "cherry"]`),
			expected: JsonArray{
				"apple", "banana", "cherry",
			},
			shouldFail: false,
		},
		{
			input: []byte(`[true, null, 42, "text"]`),
			expected: JsonArray{
				true, nil, int64(42), "text",
			},
			shouldFail: false,
		},
		{
			input: []byte(`{"outer": {"inner": "value"}}`),
			expected: JsonObject{
				"outer": JsonObject{
					"inner": "value",
				},
			},
			shouldFail: false,
		},
		{
			input:      []byte(`[]`),
			expected:   JsonArray{},
			shouldFail: false,
		},
		{
			input:      []byte(`{}`),
			expected:   JsonObject{},
			shouldFail: false,
		},
		{
			input:      []byte(`{"key": "value"`), // Invalid JSON, missing closing brace
			expected:   nil,
			shouldFail: true,
		},
		{
			input:      []byte(`{"key": 012}`), // Invalid JSON with malformed number
			expected:   nil,
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tc.input)
			if err != nil {
				t.Fatalf("Failed to tokenize input: %v", err)
			}

			result, err := Parse(tokens)
			if tc.shouldFail {
				if err == nil {
					t.Fatalf("Expected an error for input %s, but got none.", tc.input)
				}
				return // Skip further checks if we expect failure
			}

			if err != nil {
				t.Fatalf("Parse(%s) returned an error: %v", tc.input, err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Parse(%s) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
