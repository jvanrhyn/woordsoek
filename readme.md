# Woordsoek

Woordsoek is a Text User Interface (TUI) for searching words in a text file based on specific criteria. It supports loading environment variables from a `.env` file.

## Features

- **Word Search**: Search for words that contain a specific single letter and are composed of characters from a given 6-character string.
- **Length Filtering**: Filter words based on a specified length.
- **Language selection**: Choose a dictionary by language (e.g., `af-za`, `en`) via environment variable, CLI flag, or TUI input.

## Usage

```bash
go run main
```

The tool operates interactively, prompting the user for input values:
- **Single Letter**: A single letter that must be present in the words.
- **6-Character String**: A string of 6 characters that the words can be composed of.
- **Word Length**: (Optional) The exact length of the words to search for.

Non-interactive (CLI) mode is available for scripting and automation:

```bash
go run main -single a -chars lph -lang en -length 0
```

- `-single`: required single character that must appear in a matching word
- `-chars`: additional allowed characters (the single letter is included automatically)
- `-length`: exact length to match (use `0` for length-agnostic and default min-length behavior)
- `-lang`: language code to choose which dictionary file in `dictionaries/` to use (overrides `WBLANG` env var)

## Environment Variables

The tool loads environment variables from a `.env` file. The primary variable used is:

- **`WBLANG`**: Specifies the language dictionary to use (default is `af-za`).

When both `WBLANG` and the `-lang` flag are provided, the CLI `-lang` flag takes precedence. In the TUI, the language input (fourth input) overrides `WBLANG` if set.

## How It Works

1. **Initialization**: Loads environment variables from a `.env` file.
2. **Command Parsing**: Parses user inputs to determine the operation mode.
3. **Word Search**: Searches for matching words in the specified dictionary file.

## Dictionary Files

The tool uses dictionary files located in the `dictionaries/` directory. The language is specified by the `WBLANG` environment variable.

Files are named like `dictionaries/<lang>.txt` (for example `dictionaries/af-za.txt` or `dictionaries/en.txt`).

Note: dictionaries are cached in-memory on first load for performance; call `ClearDictionaryCache` (or restart the program) to pick up file changes, or use the TUI/CLI to change the selected language.

## Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.

## Testing & Benchmarks

Run the test suite:

```bash
go test ./...
```

Run fuzzers locally (may take time):

```bash
go test ./internal/woordsoek -fuzz FuzzIsValidWord -run TestNone
```

Run the microbenchmarks:

```bash
go test -bench BenchmarkSearch_ -benchmem ./internal/woordsoek
```

## Changelog (high level)

- Added input validation and Unicode-aware length checks
- Added non-interactive CLI flags (`-single`, `-chars`, `-length`, `-lang`)
- Added TUI language input (overrides `WBLANG`)
- Improved performance: dictionary caching and parallel search for large dictionaries
- Added fuzz tests, edge-case tests, and microbenchmarks


## License

This project is licensed under the MIT License.
