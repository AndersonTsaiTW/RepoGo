# RepoGo

A CLI tool that helps you summarize your repository information. You can use the generated summary as input for LLM models, documentation, or other analysis tasks.

## Features

- 🔍 **Smart File Scanning**: Automatically scans project files and generates structured summaries
- 📊 **Multiple Output Formats**: Supports Markdown and JSON output formats
- 🎯 **Flexible Filtering**: Supports include/exclude glob patterns for file filtering
- 📏 **Token Estimation**: Can display estimated token counts for LLM input control
- 🚀 **Efficient Processing**: Automatically handles binary files with file size limits
- 📋 **Git Integration**: Automatically retrieves Git repository information (branch, commits, etc.)

## Installation

### Build from Source

```bash
git clone https://github.com/AndersonTsaiTW/RepoGo.git
cd RepoGo
make build
```

The compiled binary will be located at `bin/repogo`

### Cross-platform Build

```bash
make cross
```

This will compile binaries for multiple platforms, with results stored in the `dist/` directory.

## Usage

### Basic Usage

```bash
# Scan current directory
./bin/repogo

# Scan specific directory
./bin/repogo /path/to/project

# Specify output file
./bin/repogo -o summary.md

# Output JSON format
./bin/repogo -format json -o summary.json
```

### Advanced Options

```bash
# Include only specific file types
./bin/repogo -include "*.go,*.md"

# Exclude specific files
./bin/repogo -exclude "*.test,vendor/*"

# Show token estimation
./bin/repogo -tokens

# Limit maximum token count
./bin/repogo -max-tokens 50000

# Limit individual file size (bytes)
./bin/repogo -max-file-size 8192
```

## Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `-o` | Output file path | stdout |
| `-format` | Output format (markdown/json) | markdown |
| `-include` | Include file patterns (comma-separated) | All files |
| `-exclude` | Exclude file patterns (comma-separated) | None |
| `-tokens` | Show estimated token count | false |
| `-max-tokens` | Maximum token limit | 0 (unlimited) |
| `-max-file-size` | Maximum file size (bytes) | 16384 |
| `-v` | Show version | - |
| `-h` | Show help | - |

## Output Examples

### Markdown Format

Generated summaries include:

- Git repository information (branch, commits, author, etc.)
- Project directory structure
- File contents and detailed information
- Token count statistics

### JSON Format

Structured JSON output, suitable for programmatic processing and API integration.

## Use Cases

- 📝 Prepare code context for LLMs
- 📚 Generate project documentation
- 🔍 Code review and analysis
- 📊 Project structure visualization
- 🤖 Automated documentation generation

## Development

### Local Development

```bash
# Install dependencies
go mod download

# Run
go run main.go

# Test
go test ./...

# Build
make build
```

## License

See the [LICENSE](LICENSE) file for details.

## Contributing

Issues and Pull Requests are welcome!
