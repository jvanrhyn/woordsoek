package main

import (
	"log/slog"

	"github.com/jvanrhyn/woordsoek/internal/api" // Import the new api package
	configure "github.com/jvanrhyn/woordsoek/internal/config"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

func main() {

	logger := configure.SetupLogging()
	slog.SetDefault(logger)

	slog.Info("Starting Woordsoek API server")
	woordsoek.LoadVowelForms() // Initialize vowel forms
	slog.Info("Vowel forms loaded")
	api.StartAPIServer() // Start the API server
}
