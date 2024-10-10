package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
)

func startAPIServer() {
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
		return c.JSON(results)
	})

	// Start the server
	app.Listen(":3000")
}
