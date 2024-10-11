package api

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jvanrhyn/woordsoek/internal/errors"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

type SearchResponse struct {
	Parameters map[string]string `json:"parameters"`
	Count      int               `json:"count"`
	Results    []string          `json:"results"`
}

func StartAPIServer() {
	slog.Info("Starting Woordsoek API server")
	app := fiber.New()

	// Define the search endpoint
	app.Get("/search", func(c *fiber.Ctx) error {
		// Extract query parameters
		locale := c.Get("x-locale")       // Get the x-locale header
		filename := "dictionaries/en.txt" // Default filename

		// Check if the x-locale header is present
		if locale != "" {
			filename = "dictionaries/" + locale + ".txt" // Set filename based on header
		}

		slog.Info("Searching for words in " + filename)

		singleLetter := c.Query("singleLetter")
		sixCharString := c.Query("sixCharString")
		lengthStr := c.Query("length")

		// Convert length to int
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			length = 0 // Default value if conversion fails
		}

		// Call the search function from woordsoek package
		results, err := woordsoek.SearchForMatchingWords(filename, singleLetter, sixCharString, length)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errors.CustomError{Message: err.Error()})
		}

		// Create the response object
		response := SearchResponse{
			Parameters: map[string]string{
				"singleLetter":  singleLetter,
				"sixCharString": sixCharString,
				"length":        lengthStr,
			},
			Count:   len(results),
			Results: results,
		}

		slog.Info("Found", "wordcount", len(results))
		return c.JSON(response)
	})

	// Start the server
	app.Listen(":3000")
}
