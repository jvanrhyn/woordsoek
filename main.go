package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/jvanrhyn/woordsoek/internal"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

func main() {
	woordsoek.LoadVowelForms() // Initialize vowel forms

	flags := tui.Flags{
		Length: 0,
	}

	p := tea.NewProgram(tui.InitializeModel(flags), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
