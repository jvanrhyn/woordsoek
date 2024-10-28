package woordsoek

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIsValidWord(t *testing.T) {
	tests := []struct {
		word          string
		singleLetter  string
		sixCharString string
		expected      bool
	}{
		{"hello", "h", "ello", true},
		{"world", "w", "orld", true},
		{"test", "t", "est", true},
		{"invalid", "x", "yz", false},
		{"", "a", "b", false},
	}

	for _, test := range tests {
		result := IsValidWord(test.word, test.singleLetter, test.sixCharString)
		if result != test.expected {
			t.Errorf("IsValidWord(%q, %q, %q) = %v; expected %v", test.word, test.singleLetter, test.sixCharString, result, test.expected)
		}
	}
}

func TestSearchForMatchingWords(t *testing.T) {
	LoadVowelForms() // Ensure vowel forms are loaded

	// Create a temporary file for testing
	tempFile := "test_words.txt"
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile)

	// Write test words to the temporary file
	words := []string{"hello", "world", "test", "word", "example"}
	if err := os.WriteFile(tempFile, []byte(strings.Join(words, "\n")), 0644); err != nil {
		t.Fatalf("Failed to write test words to file: %v", err)
	}

	tests := []struct {
		filename      string
		singleLetter  string
		sixCharString string
		length        int
		expected      []string
	}{
		{tempFile, "h", "ello", 0, []string{"hello"}},
		{tempFile, "w", "orld", 0, []string{"word", "world"}},
		{tempFile, "e", "xample", 0, []string{"example"}},
		{tempFile, "t", "est", 4, []string{"test"}},
		{tempFile, "x", "", 0, []string{}}, // No matches
	}

	for _, test := range tests {
		result, err := SearchForMatchingWords(test.filename, test.singleLetter, test.sixCharString, test.length)
		if err != nil {
			t.Errorf("SearchForMatchingWords(%q, %q, %q, %d) returned an error: %v", test.filename, test.singleLetter, test.sixCharString, test.length, err)
			continue
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SearchForMatchingWords(%q, %q, %q, %d) = %v; expected %v", test.filename, test.singleLetter, test.sixCharString, test.length, result, test.expected)
		}
	}
}
