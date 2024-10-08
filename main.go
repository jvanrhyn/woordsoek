package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

type (
	VowelForms map[rune]string
	Flags      struct {
		CleanupFlag   bool
		Filename      string
		OutputFile    string
		SingleLetter  string
		SixCharString string
		Length        int
	}
)

var (
	vowelForms VowelForms
)

func init() {
	loadEnv()
}

func main() {
	flags := Flags{}

	flag.BoolVar(&flags.CleanupFlag, "cleanup", false, "Run cleanup on a given file (alias: -c)")
	flag.BoolVar(&flags.CleanupFlag, "cl", false, "Run cleanup on a given file")

	flag.StringVar(&flags.Filename, "file", "", "The input file (required for cleanup) (alias: -f)")
	flag.StringVar(&flags.Filename, "f", "", "The input file (required for cleanup)")

	flag.StringVar(&flags.OutputFile, "output", "", "The output file (required for cleanup) (alias: -o)")
	flag.StringVar(&flags.OutputFile, "o", "", "The output file (required for cleanup)")

	flag.StringVar(&flags.SingleLetter, "letter", "", "A single letter to search for (alias: -l)")
	flag.StringVar(&flags.SingleLetter, "l", "", "A single letter to search for")

	flag.StringVar(&flags.SixCharString, "chars", "", "A 6-character string (alias: -c)")
	flag.StringVar(&flags.SixCharString, "c", "", "A 6-character string")

	flag.IntVar(&flags.Length, "length", 0, "Optional length of words to match (alias: -len)")
	flag.IntVar(&flags.Length, "len", 0, "Optional length of words to match")

	flag.Parse()

	if flags.CleanupFlag {
		if flags.Filename == "" || flags.OutputFile == "" {
			fmt.Println("Usage: -cleanup -file <filename> -output <outputFile>")
			return
		}
		fmt.Println("--- Cleanup called ---")
		cleanUp(flags.Filename, flags.OutputFile)
		return
	}

	if flags.SingleLetter == "" || flags.SixCharString == "" {
		fmt.Println("Usage: -letter <single-letter> -chars <6-char-string> [-length <length>]")
		return
	}

	if len(flags.SingleLetter) != 1 {
		fmt.Println("The letter parameter must be a single letter.")
		return
	}

	if len(flags.SixCharString) != 6 {
		fmt.Println("The chars parameter must be a string of 6 characters.")
		return
	}

	lang := os.Getenv("WBLANG")
	if lang == "" {
		lang = "af-za"
	}
	filenamePath := filepath.Join("dictionaries", lang+".txt")

	// Define a map of vowels to their different forms
	vowelForms = VowelForms{
		'a': "àáâãäå",
		'e': "èéêë",
		'i': "ìíîï",
		'o': "òóôõö",
		'u': "ùúûü",
	}

	// Update sixCharString to include different forms of vowels
	updatedSixCharString := flags.SixCharString
	for _, char := range flags.SixCharString {
		if forms, exists := vowelForms[char]; exists {
			updatedSixCharString += forms
		}
	}
	flags.SixCharString = updatedSixCharString

	searchForMatchingWords(filenamePath, flags.SingleLetter, flags.SixCharString, flags.Length)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func searchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var results []string

	for scanner.Scan() {
		word := scanner.Text()
		if strings.Contains(word, singleLetter) && isValidWord(word, singleLetter, sixCharString) {
			if length == 0 || len(word) == length {
				results = append(results, word)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
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

	fmt.Println("Matching words:", results)

	fmt.Println("Number of matching words:", len(results))
}

func isValidWord(word, singleLetter, sixCharString string) bool {
	allowedChars := strings.ToLower(singleLetter + sixCharString)
	word = strings.ToLower(word)
	for _, char := range word {
		if !strings.ContainsRune(allowedChars, char) {
			return false
		}
	}
	return true
}

func cleanUp(filename string, outputFileName string) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new file to write the cleaned content
	outputFile, err := os.Create("dictionaries/" + outputFileName + ".txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		// Find the index of the slash
		if index := strings.Index(line, "/"); index != -1 {
			// Remove the slash and everything after it
			line = line[:index]
		}

		// Remove numbers from the line
		line = strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return -1
			}
			return r
		}, line)

		// Write the cleaned line to the output file
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()
}
