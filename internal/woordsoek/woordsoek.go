package woordsoek

import (
	"bufio"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/jvanrhyn/woordsoek/internal/errors"
)

type VowelForms map[rune]string

var (
	vowelForms VowelForms
)

func SearchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) ([]string, error) {
	// Open the file for reading
	slog.Info("Opening file: " + filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, &errors.CustomError{Message: "Error opening file: " + err.Error()}
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)

	var results []string

	for scanner.Scan() {
		word := scanner.Text()
		slog.Info("Read word: " + word) // Debug log for each word read
		if strings.Contains(word, singleLetter) {
			slog.Info("Checking word: " + word) // Debug log before checking validity
			if IsValidWord(word, singleLetter, sixCharString) {
				slog.Info("Valid word: " + word) // Debug log for valid words
				if length == 0 || len(word) == length {
					results = append(results, word)
					slog.Info("Added word: " + word) // Debug log for added words
				}
			} else {
				slog.Info("Invalid word: " + word) // Debug log for invalid words
			}
		} else {
			slog.Info("Does not contain single letter: " + word) // Debug log for words not containing the letter
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, &errors.CustomError{Message: "Error reading file: " + err.Error()}
	}

	// Filter results for words 4 letters and longer only if length is not 0
	var filteredResults []string
	for _, word := range results {
		if length != 0 && len(word) < 4 {
			continue
		}
		filteredResults = append(filteredResults, word)
	}

	// Replace vowel forms with the base vowel in the results
	for i, word := range filteredResults {
		for vowel, forms := range vowelForms {
			for _, form := range forms {
				word = strings.ReplaceAll(word, string(form), string(vowel))
			}
		}
		filteredResults[i] = word
	}

	// Remove duplicate words from the results
	uniqueResults := make(map[string]struct{})
	for _, word := range filteredResults {
		uniqueResults[word] = struct{}{}
	}

	results = make([]string, 0, len(uniqueResults))
	for word := range uniqueResults {
		results = append(results, word)
	}
	// Sort the results alphabetically
	sort.Strings(results)

	return results, nil
}

func IsValidWord(word, singleLetter, sixCharString string) bool {
	if word == "" { // Check for empty string
		return false
	}
	allowedChars := strings.ToLower(singleLetter + sixCharString)
	slog.Info("Allowed characters: " + allowedChars) // Debug log for allowed characters
	word = strings.ToLower(word)
	for _, char := range word {
		slog.Info("Checking character: " + string(char)) // Debug log for each character checked
		if !strings.ContainsRune(allowedChars, char) {
			return false
		}
	}
	return true
}

func LoadVowelForms() {
	// Define a map of vowels to their different forms
	vowelForms = VowelForms{
		'a': "àáâãäå",
		'e': "èéêë",
		'i': "ìíîï",
		'o': "òóôõö",
		'u': "ùúûü",
	}
}
