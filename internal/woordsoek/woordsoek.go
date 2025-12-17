// Package woordsoek is used to lookup words based containing certain letters
package woordsoek

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	internalerrors "github.com/jvanrhyn/woordsoek/internal/errors"
)

type VowelForms map[rune]string

var (
	vowelForms    VowelForms
	vowelReplacer *strings.Replacer
	vowelOnce     sync.Once
)

const minDefaultLength = 4

// SearchForMatchingWords opens the dictionary file and searches for words
// that match the provided criteria. singleLetter must be exactly one rune.
func SearchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, &internalerrors.CustomError{Message: fmt.Sprintf("error opening file %s: %v", filename, err)}
	}
	defer func() { _ = file.Close() }()

	return searchFromReader(file, singleLetter, sixCharString, length)
}

// searchFromReader is the core search routine and is split out to ease testing.
func searchFromReader(r io.Reader, singleLetter, sixCharString string, length int) ([]string, error) {
	// validate singleLetter
	if singleLetter == "" {
		return nil, &internalerrors.CustomError{Message: "singleLetter must be provided"}
	}
	singleRunes := []rune(strings.ToLower(singleLetter))
	if len(singleRunes) != 1 {
		return nil, &internalerrors.CustomError{Message: "singleLetter must be a single character"}
	}

	// build allowed set for O(1) lookups
	allowed := make(map[rune]struct{})
	for _, r := range strings.ToLower(singleLetter + sixCharString) {
		allowed[r] = struct{}{}
	}

	scanner := bufio.NewScanner(r)
	var candidates []string
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		// skip quickly if it doesn't contain the requested rune
		if !strings.ContainsRune(strings.ToLower(word), singleRunes[0]) {
			continue
		}
		if isValidWordWithAllowedSet(word, allowed) {
			// length check: if length is specified, require exact match; otherwise require minDefaultLength
			runeLen := utf8.RuneCountInString(word)
			if length > 0 {
				if runeLen != length {
					continue
				}
			} else {
				if runeLen < minDefaultLength {
					continue
				}
			}
			candidates = append(candidates, normalizeWord(word))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, &internalerrors.CustomError{Message: fmt.Sprintf("error reading input: %v", err)}
	}

	// unique and sort
	uniq := make(map[string]struct{})
	for _, w := range candidates {
		uniq[w] = struct{}{}
	}
	results := make([]string, 0, len(uniq))
	for w := range uniq {
		results = append(results, w)
	}
	sort.Strings(results)
	slog.Debug("Search completed", "matches", len(results))
	return results, nil
}

// IsValidWord kept for compatibility with tests: it builds the allowed set and delegates.
func IsValidWord(word, singleLetter, sixCharString string) bool {
	if word == "" {
		return false
	}
	allowed := make(map[rune]struct{})
	for _, r := range strings.ToLower(singleLetter + sixCharString) {
		allowed[r] = struct{}{}
	}
	return isValidWordWithAllowedSet(word, allowed)
}

func isValidWordWithAllowedSet(word string, allowed map[rune]struct{}) bool {
	for _, r := range strings.ToLower(word) {
		if _, ok := allowed[r]; !ok {
			return false
		}
	}
	return true
}

func normalizeWord(word string) string {
	vowelOnce.Do(func() {
		// Build replacer pairs from vowelForms
		var pairs []string
		for base, forms := range vowelForms {
			baseStr := string(base)
			for _, r := range forms {
				pairs = append(pairs, string(r), baseStr)
			}
		}
		if len(pairs) > 0 {
			vowelReplacer = strings.NewReplacer(pairs...)
		}
	})
	if vowelReplacer != nil {
		return strings.ToLower(vowelReplacer.Replace(word))
	}
	return strings.ToLower(word)
}

func LoadVowelForms() {
	vowelForms = VowelForms{
		'a': "àáâãäå",
		'e': "èéêë",
		'i': "ìíîï",
		'o': "òóôõö",
		'u': "ùúûü",
	}
	// reset replacer so it will be rebuilt on next normalize
	vowelOnce = sync.Once{}
	vowelReplacer = nil
}
