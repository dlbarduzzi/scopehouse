package inflector

import (
	"fmt"
	"testing"
)

func TestCapitalize(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"", ""},
		{" ", " "},
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"HELLO", "HELLO"},
		{"hello world", "Hello world"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.value), func(t *testing.T) {
			result := Capitalize(tc.value)
			if result != tc.expected {
				t.Fatalf("expected result to be %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestFormatSentence(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"", ""},
		{"  ", ""},
		{".", "."},
		{"!", "!"},
		{"?", "?"},
		{"hello", "Hello."},
		{"Hello", "Hello."},
		{"hello!", "Hello!"},
		{"hello world", "Hello world."},
		{"hello world.", "Hello world."},
		{"hello world!", "Hello world!"},
		{"hello world?", "Hello world?"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%s", i, tc.value), func(t *testing.T) {
			result := FormatSentence(tc.value)
			if result != tc.expected {
				t.Fatalf("expected result to be %s, got %s", tc.expected, result)
			}
		})
	}
}
