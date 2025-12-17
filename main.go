package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"

	tea "github.com/charmbracelet/bubbletea"
	configure "github.com/jvanrhyn/woordsoek/internal/config" // Correct import statement
	"github.com/jvanrhyn/woordsoek/internal/tui"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

func main() {
	logger := configure.SetupLogging()
	slog.SetDefault(logger)
	woordsoek.LoadVowelForms()

	// CLI flags for non-interactive usage
	var single string
	var chars string
	var length int
	var lang string
	flag.StringVar(&single, "single", "", "Single letter to include in words")
	flag.StringVar(&chars, "chars", "", "Additional allowed characters")
	flag.IntVar(&length, "length", 0, "Word length to match (0 for any)")
	flag.StringVar(&lang, "lang", "", "Language code for dictionary (overrides WBLANG env)")
	flag.Parse()

	if single != "" {
		// Non-interactive mode: search and print results
		if lang == "" {
			lang = os.Getenv("WBLANG")
		}
		if lang == "" {
			lang = "af-za"
		}
		filename := filepath.Join("dictionaries", lang+".txt")
		results, err := woordsoek.SearchForMatchingWords(filename, single, chars, length)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, w := range results {
			fmt.Println(w)
		}
		return
	}

	flags := tui.Flags{Length: 0}
	p := tea.NewProgram(tui.InitializeModel(flags), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
