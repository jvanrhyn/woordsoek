package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

type SearchResponse struct {
	Parameters map[string]string `json:"parameters"`
	Count      int               `json:"count"`
	Results    []string          `json:"results"`
}

func StartAPIServer() {
	app := fiber.New()

	// Define the search endpoint
	app.Get("/search", func(c *fiber.Ctx) error {
		// Extract query parameters
		filename := "dictionaries/en.txt" // Example filename, adjust as necessary
		singleLetter := c.Query("singleLetter")
		sixCharString := c.Query("sixCharString")
		lengthStr := c.Query("length")

		// Convert length to int
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			length = 0 // Default value if conversion fails
		}

		// Call the search function from woordsoek package
		results := woordsoek.SearchForMatchingWords(filename, singleLetter, sixCharString, length)

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

		return c.JSON(response)
	})

	// Start the server
	app.Listen(":3000")
}
