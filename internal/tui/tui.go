package tui

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	configure "github.com/jvanrhyn/woordsoek/internal/config"
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

type wordItem string

func init() {
	logger := configure.SetupLogging()
	slog.SetDefault(logger)
}

func (w wordItem) FilterValue() string {
	return string(w)
}

func (w wordItem) Title() string {
	return string(w)
}

type Model struct {
	flags        Flags
	results      []string
	loading      bool
	errorMessage string
	inputs       []textinput.Model
	focusedInput int
	currentState state
	list         list.Model
	paginator    paginator.Model
}

func InitializeModel(flags Flags) Model {
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

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10 // Default value, will be updated based on screen height

	return Model{
		flags:        flags,
		loading:      false,
		inputs:       inputs,
		focusedInput: 0,
		currentState: inputSingleLetter,
		list:         list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		paginator:    p,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "tab":
			if m.currentState == done {
				// Clear input fields
				for i := range m.inputs {
					m.inputs[i].SetValue("") // Clear the input field
				}
				m.currentState = inputSingleLetter // Reset to input state
				m.focusedInput = 0
				return m, nil
			}
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

				slog.Info("Input values",
					"SingleLetter", m.flags.SingleLetter,
					"SixCharString", m.flags.SixCharString,
					"Length", m.flags.Length,
				)

				m.currentState++
				if m.focusedInput < len(m.inputs)-1 {
					m.focusedInput++
				}
				m.inputs[m.focusedInput].Focus()
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		// Update the paginator's PerPage based on the terminal height
		m.paginator.PerPage = msg.Height / 2 // Adjust this value based on your layout
	}

	// Only update the current input field
	if m.currentState <= inputLength {
		m.inputs[m.focusedInput], _ = m.inputs[m.focusedInput].Update(msg)
	}

	// Update the paginator model if in done state
	if m.currentState == done {
		var cmd tea.Cmd
		m.paginator, cmd = m.paginator.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) searchWords() Model {
	m.loading = true // Set loading to true when starting the search
	lang := os.Getenv("WBLANG")
	slog.Info("Language from environment", "lang", lang)
	if lang == "" {
		lang = "af-za"
		slog.Info("Language set to", "lang", lang)
	}
	filenamePath := filepath.Join("dictionaries", lang+".txt")

	woordsoek.LoadVowelForms()
	var err error
	m.results, err = woordsoek.SearchForMatchingWords(filenamePath, m.flags.SingleLetter, m.flags.SixCharString, m.flags.Length)
	if err != nil {
		m.errorMessage = "Error searching for words: " + err.Error()
		slog.Error("Error searching for words", "error", err)
	}
	m.loading = false

	// Populate the list with results
	var items []list.Item
	for _, word := range m.results {
		items = append(items, wordItem(word))
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, len(items))
	m.paginator.SetTotalPages(len(items))

	return m
}

func (m Model) View() string {
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

		// Render the current page of results using the paginator
		start, end := m.paginator.GetSliceBounds(len(m.results))
		var b strings.Builder
		b.WriteString("\nMatching Words:\n\n")
		for _, item := range m.results[start:end] {
			b.WriteString("  • " + item + "\n")
		}
		b.WriteString("\n" + m.paginator.View())
		b.WriteString("\n\nPress 'esc' to quit. Press 'tab' to restart.\n")
		return b.String()
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
