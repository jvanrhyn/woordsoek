package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type (
	VowelForms map[rune]string
	Flags      struct {
		CleanupFlag   bool
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

type state int

const (
	inputSingleLetter state = iota
	inputSixCharString
	inputLength
	done
)

type model struct {
	flags        Flags
	results      []string
	loading      bool
	errorMessage string
	inputs       []textinput.Model
	focusedInput int
	currentState state
}

func main() {
	flags := Flags{}

	p := tea.NewProgram(initializeModel(flags), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initializeModel(flags Flags) model {
	inputs := make([]textinput.Model, 3)

	// SingleLetter input
	input := textinput.New()
	input.Placeholder = "Single Letter"
	input.Focus()
	inputs[0] = input

	// SixCharString input
	input = textinput.New()
	input.Placeholder = "6-Character String"
	inputs[1] = input

	// Length input
	input = textinput.New()
	input.Placeholder = "Word Length (0 for any)"
	inputs[2] = input

	return model{
		flags:        flags,
		loading:      false,
		inputs:       inputs,
		focusedInput: 0,
		currentState: inputSingleLetter,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "tab":
			return initializeModel(m.flags), nil
		case "enter":
			if m.currentState <= inputLength {
				if m.currentState == inputSingleLetter {
					m.flags.SingleLetter = m.inputs[0].Value()
				} else if m.currentState == inputSixCharString {
					m.flags.SixCharString = m.inputs[1].Value()
				} else if m.currentState == inputLength {
					length, err := strconv.Atoi(m.inputs[2].Value())
					if err != nil {
						m.errorMessage = "Length must be a number"
						return m, nil
					}
					m.flags.Length = length
					m.currentState = done
					return m.searchWords(), nil
				}

				m.currentState++
				m.focusedInput++
				m.inputs[m.focusedInput].Focus()
				return m, nil
			}
		}
	}

	// Only update the current input field
	if m.currentState <= inputLength {
		m.inputs[m.focusedInput], _ = m.inputs[m.focusedInput].Update(msg)
	}

	return m, nil
}

func (m model) searchWords() model {
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
	updatedSixCharString := m.flags.SixCharString
	for _, char := range m.flags.SixCharString {
		if forms, exists := vowelForms[char]; exists {
			updatedSixCharString += forms
		}
	}
	m.flags.SixCharString = updatedSixCharString

	m.results = searchForMatchingWords(filenamePath, m.flags.SingleLetter, m.flags.SixCharString, m.flags.Length)
	m.loading = false

	return m
}

func (m model) View() string {
	if m.loading {
		return "Loading..."
	}

	if m.errorMessage != "" {
		return fmt.Sprintf("Error: %s\nPress 'esc' to quit.", m.errorMessage)
	}

	if m.currentState == done {
		if len(m.results) == 0 {
			return "No matching words found.\nPress 'esc' to quit."
		}

		resultStr := "Matching words:\n"
		for _, word := range m.results {
			resultStr += word + "\n"
		}
		resultStr += fmt.Sprintf("\nNumber of matching words: %d\n", len(m.results))
		resultStr += "\nPress 'esc' to quit. Press 'tab' to restart."

		return resultStr
	}

	var b strings.Builder
	b.WriteString("Input Values (Press 'Enter' to continue):\n\n")

	for i := range m.inputs {
		if m.currentState <= inputLength && i > 2 {
			break
		}
		b.WriteString(m.inputs[i].View())
		if i == m.focusedInput {
			b.WriteString(" ←")
		}
		b.WriteString("\n")
	}

	b.WriteString("\nPress 'esc' to quit, 'tab' to restart.")

	return b.String()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func searchForMatchingWords(filename string, singleLetter string, sixCharString string, length int) []string {
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
		if strings.Contains(word, singleLetter) && isValidWord(word, singleLetter, sixCharString) {
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
