# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`envee` is a Go CLI tool that injects environment variables from a file into a child process. Usage: `envee [-f .env] -- command [args...]`

## Commands

```bash
# Build for current platform
go build -o envee .

# Run tests
go test ./...

# Run a single test
go test -run TestName .

# Build release archives for all platforms (linux/darwin × amd64/arm64)
./build.sh

# Format code (auto-runs on save via VS Code config)
gofmt -w .

# Vet
go vet ./...
```

## Architecture

The project is two files at the root:

- **`main.go`**: CLI entrypoint. Uses `flag` to parse `-f` (env file path, default `.env`). Remaining args after `--` are the command to execute. Merges env file vars over `os.Environ()`, then runs the command via `exec.Command` with stdin/stdout/stderr passed through. Propagates the child's exit code.
- **`utils.go`**: Two functions — `parseEnvFile()` parses `KEY=VALUE` lines (supports comments, `export` prefix, quote stripping; no variable expansion) and `mergeEnv()` overlays the parsed map onto the base environment slice.

## env file format

- `KEY=VALUE` pairs, one per line
- Lines starting with `#` are ignored
- Optional `export KEY=VALUE` syntax supported
- Surrounding single or double quotes on values are stripped
- No variable expansion, no escape sequences, no multi-line values
