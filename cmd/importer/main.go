package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Database connection parameters
	connStr := "user=johanvanrhyn dbname=woordsoek sslmode=disable" // Update with your username
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Iterate over files in the dictionaries folder
	err = filepath.Walk("dictionaries", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Get the locale from the filename
			locale := info.Name()
			locale = locale[:len(locale)-len(filepath.Ext(locale))]
			fmt.Println("Importing words for locale:", locale)

			// Split content into words (assuming each word is separated by new lines)
			words := strings.Split(string(content), "\n")
			insertCount := 0 // Counter for inserts
			for _, word := range words {
				// Insert into the database
				_, err := db.Exec("INSERT INTO words (word, locale, in_use) VALUES ($1, $2, $3)", word, locale, true)
				if err != nil {
					return err
				}
				insertCount++ // Increment the counter

				// Report progress every 500 inserts
				if insertCount%500 == 0 {
					fmt.Printf("Inserted %d words for locale: %s\n", insertCount, locale)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Words imported successfully.")
}
