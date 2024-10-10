package woordsoek

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type VowelForms map[rune]string

var (
	vowelForms VowelForms
)

func SearchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) []string {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return []string{"Error opening file: " + err.Error()}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var results []string

	for scanner.Scan() {
		word := scanner.Text()
		if strings.Contains(word, singleLetter) && IsValidWord(word, singleLetter, sixCharString) {
			if length == 0 || len(word) == length {
				results = append(results, word)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return []string{"Error reading file: " + err.Error()}
	}

	// Filter results for words 4 letters and longer
	var filteredResults []string
	for _, word := range results {
		if len(word) >= 4 {
			filteredResults = append(filteredResults, word)
		}
	}

	results = filteredResults

	// Replace vowel forms with the base vowel in the results
	for i, word := range results {
		for vowel, forms := range vowelForms {
			for _, form := range forms {
				word = strings.ReplaceAll(word, string(form), string(vowel))
			}
		}
		results[i] = word
	}

	// Remove duplicate words from the results
	uniqueResults := make(map[string]struct{})
	for _, word := range results {
		uniqueResults[word] = struct{}{}
	}

	results = make([]string, 0, len(uniqueResults))
	for word := range uniqueResults {
		results = append(results, word)
	}
	// Sort the results alphabetically
	sort.Strings(results)

	return results
}

func IsValidWord(word, singleLetter, sixCharString string) bool {
	allowedChars := strings.ToLower(singleLetter + sixCharString)
	word = strings.ToLower(word)
	for _, char := range word {
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
