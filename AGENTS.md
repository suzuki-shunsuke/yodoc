# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

yodoc is a CLI tool that maintains documents containing commands and their execution results. It executes commands from template files, tests their results, and generates documents. This solves the problem of keeping documentation with embedded command outputs up to date.

## Build and Development Commands

```sh
# Run tests
go test ./... -race -covermode=atomic

# Run linter
go vet ./...
golangci-lint run

# Generate JSON schema
go run ./cmd/gen-jsonschema
```

## Architecture

The codebase follows a standard Go CLI structure with clear separation of concerns:

### Entry Points
- `cmd/yodoc/main.go` - CLI entry point
- `cmd/gen-jsonschema/main.go` - JSON schema generator for the config file

### Core Packages
- `pkg/cli/` - CLI command definitions using urfave/cli/v3
- `pkg/controller/run/` - Main execution logic for the `yodoc run` command
- `pkg/controller/initcmd/` - Logic for `yodoc init` command
- `pkg/parser/` - Parses template files into blocks (hidden, run, check, other, out)
- `pkg/render/` - Template rendering with sprig functions (some functions removed for security)
- `pkg/frontmatter/` - YAML front matter parsing for templates
- `pkg/config/` - Configuration file reading and validation
- `pkg/expr/` - Expression evaluation using expr-lang/expr for check assertions

### Processing Flow
1. Find and read config file (`yodoc.yaml` or `.yodoc.yaml`)
2. Discover template files (`.md`/`.mdx` in src directory)
3. Parse each template: extract front matter, then parse into blocks
4. Process blocks sequentially - execute commands, run checks, render templates
5. Write generated documents to destination directory

### Block Types
The parser (`pkg/parser/parser.go`) identifies these block types via annotations:
- `HiddenBlock` (`#-yodoc hidden`) - Execute but don't output
- `MainBlock` (`#-yodoc run` / `#!yodoc run`) - Execute and output
- `CheckBlock` (`#-yodoc check`) - Validate previous command results
- `OtherBlock` - Regular code blocks
- `OutBlock` - Content outside code blocks

### Template Context
Command results are available as template variables: `Command`, `ExitCode`, `Stdout`, `Stderr`, `CombinedOutput`, `RunError`.

### Security Considerations
- Some sprig functions are disabled (`env`, `expandenv`, `getHostByName`) for security
- Templates execute shell commands, so only run on trusted files
