package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanUp(filename string, outputFileName string) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Create a new file to write the cleaned content
	outputFile, err := os.Create("dictionaries/" + outputFileName + ".txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

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
	_ = writer.Flush()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cleaner.go <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	cleanUp(inputFile, outputFile)
}
