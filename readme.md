# Woordsoek

Woordsoek is a command-line tool for searching words in a text file based on specific criteria. It supports loading environment variables from a `.env` file and cleaning up text files by removing comments.

## Features

- **Word Search**: Search for words that contain a specific single letter and are composed of characters from a given 6-character string.
- **Length Filtering**: Filter words based on a specified length.
- **Text Cleanup**: Clean up text files by removing comments and numbers.

## Usage

```bash
main <single-letter> <6-char-string> [length]
main <filename>.txt
```

- **`<single-letter>`**: A single letter that must be present in the words.
- **`<6-char-string>`**: A string of 6 characters that the words can be composed of.
- **`[length]`**: (Optional) The exact length of the words to search for.
- **`<filename>.txt`**: The name of the text file to clean up.

## Environment Variables

The tool loads environment variables from a `.env` file. The primary variable used is:

- **`WBLANG`**: Specifies the language dictionary to use (default is `af-za`).

## How It Works

1. **Initialization**: Loads environment variables from a `.env` file.
2. **Command Parsing**: Parses command-line arguments to determine the operation mode.
3. **Word Search**: Searches for matching words in the specified dictionary file.
4. **Text Cleanup**: Cleans up text files by removing comments and numbers, and writes the cleaned content to a new file.

## Dictionary Files

The tool uses dictionary files located in the `dictionaries/` directory. The language is specified by the `WBLANG` environment variable.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.

## License

This project is licensed under the MIT License.
