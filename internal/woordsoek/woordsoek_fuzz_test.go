package woordsoek

import (
	"sort"
	"strings"
	"testing"
	"unicode/utf8"
)

// Fuzz IsValidWord to ensure it never panics and respects the allowed set
func FuzzIsValidWord(f *testing.F) {
	// seeds
	f.Add("hello", "h", "ello")
	f.Add("café", "c", "afe")
	f.Add("", "a", "b")

	f.Fuzz(func(t *testing.T, word, single, chars string) {
		// Build allowed set the same way IsValidWord does
		allowed := make(map[rune]struct{})
		for _, r := range strings.ToLower(single + chars) {
			allowed[r] = struct{}{}
		}

		// call function under test
		ok := IsValidWord(word, single, chars)

		// If ok is true, then every rune in word must be in allowed set
		if ok {
			for _, r := range strings.ToLower(word) {
				if _, exists := allowed[r]; !exists {
					t.Fatalf("IsValidWord returned true but rune %q not allowed (word=%q single=%q chars=%q)", r, word, single, chars)
				}
			}
		}

		// ensure function handles unicode input (no panic) and is deterministic
		_ = utf8.RuneCountInString(word)
	})
}

// Fuzz searchFromReader to verify invariants for returned results
func FuzzSearchFromReader(f *testing.F) {
	seedText := "hello\nworld\ncafé\nCafé\nexample\nword\n"
	f.Add(seedText, "h", "ello", 0)
	f.Add(seedText, "c", "afe", 0)
	f.Add(seedText, "x", "", 0)

	f.Fuzz(func(t *testing.T, dict, single, chars string, length int) {
		// limit length to reasonable range to avoid slow fuzz runs
		if length < 0 || length > 50 {
			return
		}

		// If single not a single rune, we expect an error
		if single == "" || utf8.RuneCountInString(single) != 1 {
			_, err := searchFromReader(strings.NewReader(dict), single, chars, length)
			if err == nil {
				t.Fatalf("expected error for invalid single runes; single=%q", single)
			}
			return
		}

		res, err := searchFromReader(strings.NewReader(dict), single, chars, length)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// verify invariants for each returned word
		singleRune := []rune(strings.ToLower(single))[0]
		allowed := make(map[rune]struct{})
		for _, r := range strings.ToLower(single + chars) {
			allowed[r] = struct{}{}
		}

		for _, w := range res {
			lw := strings.ToLower(w)
			if !strings.ContainsRune(lw, singleRune) {
				t.Fatalf("result %q does not contain required rune %q", w, singleRune)
			}
			for _, r := range lw {
				if _, ok := allowed[r]; !ok {
					t.Fatalf("result %q contains disallowed rune %q", w, r)
				}
			}
			if length > 0 {
				if utf8.RuneCountInString(w) != length {
					t.Fatalf("result %q does not meet requested length %d", w, length)
				}
			} else if utf8.RuneCountInString(w) < minDefaultLength {
				t.Fatalf("result %q is shorter than default minimum %d", w, minDefaultLength)
			}
		}

		// results must be unique and sorted
		if !sort.StringsAreSorted(res) {
			t.Fatalf("results not sorted: %v", res)
		}
		if len(res) > 1 {
			// check uniqueness
			for i := 1; i < len(res); i++ {
				if res[i] == res[i-1] {
					t.Fatalf("duplicate result found: %v", res[i])
				}
			}
		}
	})
}
