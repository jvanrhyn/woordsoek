package main

import (
	"github.com/jvanrhyn/woordsoek/internal/api" // Import the new api package
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

func main() {
	woordsoek.LoadVowelForms() // Initialize vowel forms
	api.StartAPIServer()       // Start the API server
}
