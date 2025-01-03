# Dead Links Scraper

**Dead Links Scraper** is a CLI tool built in Go that scans a website for dead links. It crawls the pages within the given domain and provides a list of dead links, along with their status codes.

## Usage

Run the CLI tool with the required options:

```bash
./deadlinks --url <website-url> [options]
```

### Options

- `-u, --url`: **Required.** The URL of the webpage to analyze.
- `-v, --verbose`: Enables verbose output for debugging.
- `-t, --threads`: Number of concurrent threads for link checking. Default is 4.

### Examples

1. Basic usage:
   ```bash
   ./deadlinks --url https://example.com
   ```

2. Using verbose mode:
   ```bash
   ./deadlinks --url https://example.com --verbose
   ```

3. Specifying the number of threads:
   ```bash
   ./deadlinks --url https://example.com --threads 8
   ```

JavaScript-rendered links are not supported as the tool does not execute JavaScript.
