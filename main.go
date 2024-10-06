// Package main provides a command-line tool for searching words in a text file
// based on specific criteria. It supports loading environment variables from
// a .env file and cleaning up text files by removing comments.
//
// The main functionality includes:
//   - Searching for words that contain a specific single letter and are composed
//     of characters from a given 6-character string.
//   - Filtering words based on a specified length.
//   - Cleaning up text files by removing comments.
//
// Usage:
//
//	main <single-letter> <6-char-string> [length]
//	main <filename>.txt
//
// Functions:
//   - init: Initializes the environment by loading variables from a .env file.
//   - main: The entry point of the application. It parses command-line arguments
//     and calls the appropriate functions based on the input.
//   - loadEnv: Loads environment variables from a .env file.
//   - searchForMatchingWords: Searches for words in a file that match the given
//     criteria and prints the results.
//   - isValidWord: Checks if a word is valid based on the allowed characters.
//   - cleanUp: Cleans up a text file by removing comments and writes the cleaned
//     content to a new file.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	loadEnv()
}

func main() {

	if len(os.Args) == 3 && (strings.HasSuffix(os.Args[1], ".txt") || strings.HasSuffix(os.Args[1], ".dic")) {
		fmt.Println("--- Cleanup called ---")
		cleanUp(os.Args[1], os.Args[2])
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: main <single-letter> <6-char-string> [length]")
		return
	}

	singleLetter := os.Args[1]
	if len(singleLetter) != 1 {
		fmt.Println("The first parameter must be a single letter.")
		return
	}

	sixCharString := os.Args[2]
	if len(sixCharString) != 6 {
		fmt.Println("The second parameter must be a string of 6 characters.")
		return
	}

	var length int
	if len(os.Args) > 3 {
		var err error
		length, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("The third parameter must be an integer.")
			return
		}
	} else {
		length = 0 // Default value if the length parameter is not provided
	}

	lang := os.Getenv("WBLANG")
	if lang == "" {
		lang = "af-za"
	}
	filename := filepath.Join("dictionaries", lang+".txt")

	searchForMatchingWords(filename, singleLetter, sixCharString, length)

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
