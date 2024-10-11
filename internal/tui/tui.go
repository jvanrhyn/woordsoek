package tui

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

type Flags struct {
	SingleLetter  string
	SixCharString string
	Length        int
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

func InitializeModel(flags Flags) model {
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
			return InitializeModel(m.flags), nil
		case "enter":
			if m.currentState <= inputLength {
				if m.currentState == inputSingleLetter {
					m.flags.SingleLetter = m.inputs[0].Value()
				} else if m.currentState == inputSixCharString {
					m.flags.SixCharString = m.inputs[1].Value()
				} else if m.currentState == inputLength {
					length, err := strconv.Atoi(m.inputs[2].Value())
					if err != nil {
						m.flags.Length = 0
					} else {
						m.flags.Length = length
					}
					m.currentState = done
					return m.searchWords(), nil
				}

				m.currentState++
				if m.focusedInput < len(m.inputs)-1 {
					m.focusedInput++
				}
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
	m.loading = true // Set loading to true when starting the search
	lang := os.Getenv("WBLANG")
	if lang == "" {
		m.errorMessage = "Language not set. Defaulting to 'af-za'."
		lang = "af-za"
	}
	filenamePath := filepath.Join("dictionaries", lang+".txt")

	woordsoek.LoadVowelForms()
	var err error
	m.results, err = woordsoek.SearchForMatchingWords(filenamePath, m.flags.SingleLetter, m.flags.SixCharString, m.flags.Length)
	if err != nil {
		m.errorMessage = "Error searching for words: " + err.Error()
	}
	m.loading = false

	return m
}

func (m model) View() string {
	if m.loading {
		return "Loading..."
	}

	if m.errorMessage != "" {
		return "Error: " + m.errorMessage + "\nPress 'esc' to quit."
	}

	if m.currentState == done {
		if len(m.results) == 0 {
			return "No matching words found.\nPress 'esc' to quit."
		}

		resultStr := "Matching words:\n"
		for _, word := range m.results {
			resultStr += word + "\n"
		}
		resultStr += "\nNumber of matching words: " + strconv.Itoa(len(m.results)) + "\n"
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
