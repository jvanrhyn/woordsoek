package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	var wg sync.WaitGroup

	// Iterate over files in the dictionaries folder
	err = filepath.Walk("dictionaries", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			wg.Add(1) // Increment the WaitGroup counter
			go func(path string, locale string) {
				defer wg.Done() // Decrement the counter when the goroutine completes

				// Read file content
				content, err := os.ReadFile(path)
				if err != nil {
					log.Printf("Error reading file %s: %v\n", path, err)
					return
				}

				fmt.Println("Importing words for locale:", locale)

				// Split content into words (assuming each word is separated by new lines)
				words := strings.Split(string(content), "\n")
				insertCount := 0        // Counter for inserts
				batchSize := 500        // Number of inserts per batch
				var batch []interface{} // Slice to hold words for batch insert

				for _, word := range words {
					batch = append(batch, word, locale, true)
					insertCount++ // Increment the counter

					// Execute batch insert when batch size is reached
					if insertCount%batchSize == 0 {
						query := "INSERT INTO words (word, locale, in_use) VALUES "
						valueStrings := make([]string, batchSize)
						for i := 0; i < batchSize; i++ {
							valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
						}
						_, err := db.Exec(query+strings.Join(valueStrings, ","), batch...)
						if err != nil {
							log.Printf("Error inserting batch for locale %s: %v\n", locale, err)
							return
						}
						fmt.Printf("Inserted %d words for locale: %s\n", insertCount, locale)
						batch = nil // Reset batch
					}
				}

				// Insert any remaining words in the batch
				if len(batch) > 0 {
					query := "INSERT INTO words (word, locale, in_use) VALUES "
					valueStrings := make([]string, len(batch)/3)
					for i := 0; i < len(batch)/3; i++ {
						valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
					}
					_, err := db.Exec(query+strings.Join(valueStrings, ","), batch...)
					if err != nil {
						log.Printf("Error inserting remaining batch for locale %s: %v\n", locale, err)
						return
					}
					fmt.Printf("Inserted %d words for locale: %s\n", insertCount, locale)
				}
			}(path, info.Name()[:len(info.Name())-len(filepath.Ext(info.Name()))]) // Pass the locale
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Words imported successfully.")
}
