# Woordsoek

Woordsoek is a Text User Interface (TUI) for searching words in a text file based on specific criteria. It supports loading environment variables from a `.env` file.

## Features

- **Word Search**: Search for words that contain a specific single letter and are composed of characters from a given 6-character string.
- **Length Filtering**: Filter words based on a specified length.

## Usage

```bash
go run main
```

The tool operates interactively, prompting the user for input values:
- **Single Letter**: A single letter that must be present in the words.
- **6-Character String**: A string of 6 characters that the words can be composed of.
- **Word Length**: (Optional) The exact length of the words to search for.

## Environment Variables

The tool loads environment variables from a `.env` file. The primary variable used is:

- **`WBLANG`**: Specifies the language dictionary to use (default is `af-za`).

## How It Works

1. **Initialization**: Loads environment variables from a `.env` file.
2. **Command Parsing**: Parses user inputs to determine the operation mode.
3. **Word Search**: Searches for matching words in the specified dictionary file.

## Dictionary Files

The tool uses dictionary files located in the `dictionaries/` directory. The language is specified by the `WBLANG` environment variable.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.

## License

This project is licensed under the MIT License.
