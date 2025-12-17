// Package woordsoek is used to lookup words based containing certain letters
package woordsoek

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
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
	dictCache     sync.Map // map[string][]string
)

const minDefaultLength = 4

// SearchForMatchingWords opens the dictionary file and searches for words
// that match the provided criteria. singleLetter must be exactly one rune.
func SearchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) ([]string, error) {
	// Use cached dictionary when available
	words, err := loadDictionary(filename)
	if err != nil {
		return nil, err
	}
	// use parallel search for larger dictionaries
	if len(words) >= 2000 {
		return searchFromWordsParallel(words, singleLetter, sixCharString, length)
	}
	return searchFromWords(words, singleLetter, sixCharString, length)
}

// loadDictionary loads a dictionary file and caches the resulting word slice.
func loadDictionary(filename string) ([]string, error) {
	if v, ok := dictCache.Load(filename); ok {
		if words, ok := v.([]string); ok {
			return words, nil
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, internalerrors.New(fmt.Sprintf("error opening file %s", filename), err)
	}
	defer func() { _ = file.Close() }()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if w := strings.TrimSpace(scanner.Text()); w != "" {
			words = append(words, w)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, internalerrors.New(fmt.Sprintf("error reading file %s", filename), err)
	}
	dictCache.Store(filename, words)
	return words, nil
}

// ClearDictionaryCache clears an entry from the dictionary cache (for tests).
func ClearDictionaryCache(filename string) {
	dictCache.Delete(filename)
}

// ClearAllDictionaryCache clears the entire dictionary cache.
func ClearAllDictionaryCache() {
	dictCache.Range(func(k, v any) bool {
		dictCache.Delete(k)
		return true
	})
}

// searchFromWords performs the search using an in-memory slice of words (serial).
func searchFromWords(words []string, singleLetter, sixCharString string, length int) ([]string, error) {
	// validate singleLetter
	if singleLetter == "" {
		return nil, internalerrors.New("singleLetter must be provided", nil)
	}
	singleRunes := []rune(strings.ToLower(singleLetter))
	if len(singleRunes) != 1 {
		return nil, internalerrors.New("singleLetter must be a single character", nil)
	}

	allowed := make(map[rune]struct{})
	for _, r := range strings.ToLower(singleLetter + sixCharString) {
		allowed[r] = struct{}{}
	}

	var candidates []string
	for _, word := range words {
		if word == "" {
			continue
		}
		if !strings.ContainsRune(strings.ToLower(word), singleRunes[0]) {
			continue
		}
		if isValidWordWithAllowedSet(word, allowed) {
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

	// unique + sort
	uniq := make(map[string]struct{}, len(candidates))
	for _, w := range candidates {
		uniq[w] = struct{}{}
	}
	results := make([]string, 0, len(uniq))
	for w := range uniq {
		results = append(results, w)
	}
	sort.Strings(results)
	return results, nil
}

// searchFromWordsParallel performs the search using multiple worker goroutines.
func searchFromWordsParallel(words []string, singleLetter, sixCharString string, length int) ([]string, error) {
	// validate inputs
	if singleLetter == "" {
		return nil, internalerrors.New("singleLetter must be provided", nil)
	}
	singleRunes := []rune(strings.ToLower(singleLetter))
	if len(singleRunes) != 1 {
		return nil, internalerrors.New("singleLetter must be a single character", nil)
	}

	allowed := make(map[rune]struct{})
	for _, r := range strings.ToLower(singleLetter + sixCharString) {
		allowed[r] = struct{}{}
	}

	// choose worker count based on CPU and input size
	maxWorkers := runtime.NumCPU()
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	numWorkers := maxWorkers
	if numWorkers > len(words) {
		numWorkers = len(words)
	}
	if numWorkers > 64 {
		numWorkers = 64
	}

	resultsCh := make(chan []string, numWorkers)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// partition the work into ranges to reduce per-item synchronization
	for i := 0; i < numWorkers; i++ {
		start := (len(words) * i) / numWorkers
		end := (len(words) * (i + 1)) / numWorkers
		go func(ws []string) {
			defer wg.Done()
			var local []string
			for _, word := range ws {
				if word == "" {
					continue
				}
				if !strings.ContainsRune(strings.ToLower(word), singleRunes[0]) {
					continue
				}
				if isValidWordWithAllowedSet(word, allowed) {
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
					local = append(local, normalizeWord(word))
				}
			}
			if len(local) > 0 {
				resultsCh <- local
			}
		}(words[start:end])
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	uniq := make(map[string]struct{})
	for slice := range resultsCh {
		for _, r := range slice {
			uniq[r] = struct{}{}
		}
	}
	results := make([]string, 0, len(uniq))
	for w := range uniq {
		results = append(results, w)
	}
	sort.Strings(results)
	return results, nil
}

// searchFromReader is the core search routine and is split out to ease testing.
func searchFromReader(r io.Reader, singleLetter, sixCharString string, length int) ([]string, error) {
	// validate singleLetter
	if singleLetter == "" {
		return nil, internalerrors.New("singleLetter must be provided", nil)
	}
	singleRunes := []rune(strings.ToLower(singleLetter))
	if len(singleRunes) != 1 {
		return nil, internalerrors.New("singleLetter must be a single character", nil)
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
		return nil, internalerrors.New("error reading input", err)
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
