# CLAUDE.md
日本語で回答して

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project called "asciinemaForLLM" that formats asciinema terminal recordings for LLM input. The project processes asciinema cast files (JSON format) and converts them into LLM-friendly structured text or CSV format. Key features include exit code extraction, working directory detection, and optimized CSV output for AI processing.

## Development Commands

### Build and Run
```bash
go build -o asciinemaForLLM
./asciinemaForLLM --help
```

### Run directly with different modes
```bash
# Format from stdin (structured text)
cat test/echo/demo.cast | go run main.go format

# Format as CSV
cat test/echo/demo.cast | go run main.go format --output=csv

# Process existing file
go run main.go file test/echo/demo.cast output.md

# Start recording session
go run main.go record my_session.cast
```

### Format code
```bash
go fmt ./...
```

### Test all cases
```bash
mise run test
```

### Individual test
```bash
cat test/echo/demo.cast | go run main.go | diff - test/echo/expectation
```

## Architecture

### Package Structure
```
/
├── main.go                    # CLI entry point with subcommand routing
├── internal/                  # Internal packages (Go standard)
│   ├── parser/
│   │   └── parser.go         # asciinema cast file parsing logic
│   ├── formatter/
│   │   └── formatter.go      # Output formatting (structured text & CSV)
│   └── cmd/
│       └── cmd.go           # CLI command implementations
└── test/                     # Test cases with demo.cast and expectation files
    ├── bash/                 # bash shell test case
    ├── echo/                 # echo command test case
    └── ls_pwd/               # ls and pwd commands test case
```

### Key Components
- **main.go**: CLI argument parsing and subcommand dispatch
- **internal/parser**: Parses .cast files, cleans terminal escape sequences, extracts commands, exit codes, and working directory
- **internal/formatter**: Outputs in structured text or CSV format with enhanced metadata
- **internal/cmd**: Implements format, record, and file subcommands with comprehensive help

## CLI Commands

### Available Subcommands
- `format` (default): Read from stdin, output to stdout
- `record [filename]`: Start asciinema recording, auto-format result
- `file <input> [output]`: Process existing .cast file

### Options
- `--output=FORMAT`: Choose output format (structured|csv)
- `--cleanup`: Remove original .cast file after processing

## Output Formats

### Structured Text (Default)
```
Terminal Session (fish shell, 148x35)
Recorded: 2025-07-08 14:14:24
Working Directory: /Users/kamonomakoto/Documents/repo/asciinemaForLLM

COMMAND: echo "Hello, world"
START TIME: 3.433s
DURATION: 2.119s
EXIT CODE: 0
OUTPUT: Hello, world
```

### CSV Format
```csv
shell,width,height,recorded,working_dir,command,start_time,duration,exit_code,output
fish,148,35,2025-07-08 14:14:24,/Users/kamonomakoto/Documents/repo/asciinemaForLLM,"echo ""Hello, world""",3.433,2.119,0,"Hello, world"
```

## File Formats

The project works with asciinema cast files (.cast) which contain:
- Header with terminal dimensions and metadata
- Array of events with timestamps and terminal output
- Format: `[timestamp, event_type, data]`

## Key Implementation Notes

- Uses Go 1.24.1 as specified in go.mod
- No external dependencies - uses only Go standard library
- Comprehensive terminal escape sequence cleaning
- Detects command execution via `cmdline_url=` markers in fish/bash output
- **Exit code extraction from `\x1b]133;D;code\x07` sequences (fish shell)**
- **Working directory extraction from OSC 7 sequences (`\x1b]7;file://`)**
- CSV output properly escapes quotes, commas, and newlines
- Backward compatible - default behavior preserved
- Module name: `github.com/Mkamono/asciinemaForLLM`

## New Features Added

### Exit Code Support
- Extracts command exit status from terminal escape sequences
- Supports both success (0) and failure (1, 127, etc.) codes
- Added `exit_code` column to CSV output
- Added `EXIT CODE` field to structured text output

### Working Directory Support
- Extracts current working directory from OSC 7 sequences
- Provides context for command execution location
- Added `working_dir` column to CSV output
- Added `Working Directory` field to structured text output

### Enhanced CSV Output
The CSV format now includes comprehensive command execution information:
```csv
shell,width,height,recorded,working_dir,command,start_time,duration,exit_code,output
```

This makes it easier for LLMs to understand:
- Command execution context (directory)
- Command success/failure status (exit code)
- Complete session metadata
