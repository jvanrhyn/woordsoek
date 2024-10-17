package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	configure "github.com/jvanrhyn/woordsoek/internal/config" // Correct import statement
	tui "github.com/jvanrhyn/woordsoek/internal/tui"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

// Comment to test signing
func main() {
	logger := configure.SetupLogging()
	slog.SetDefault(logger)
	woordsoek.LoadVowelForms()

	flags := tui.Flags{
		Length: 0,
	}

	p := tea.NewProgram(tui.InitializeModel(flags), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
