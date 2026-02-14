# envee

A lightweight Go CLI tool that injects environment variables from a `.env` file into a child process. This is useful for managing environment configuration without modifying your shell profile or cluttering your command line with environment variables.

## Installation

### Homebrew

```bash
brew tap jackjakarta/envee
brew install envee
```

### Building from Source

```bash
git clone https://github.com/jackjakarta/envee.git
cd envee
go build -o envee .
```

## Usage

```bash
envee [-f .env] -- command [args...]
```

### Options

- `-f <file>`: Path to the environment file (default: `.env`)

### Examples

Load variables from `.env` and run a command:
```bash
envee -- node server.js
```

Use a custom env file:
```bash
envee -f .env.production -- npm start
```

Pass arguments to the command:
```bash
envee -- python script.py --verbose --output=results.txt
```

## Environment File Format

The `.env` file supports the following syntax:

```bash
# Comments start with #
DATABASE_URL=postgresql://localhost/mydb
API_KEY=secret123

# Optional export prefix
export DEBUG=true

# Values with quotes (quotes are stripped)
MESSAGE="Hello World"
PATH_VAR='/usr/local/bin'

# Empty values are supported
OPTIONAL_VAR=
```

### Syntax Details

- **KEY=VALUE pairs**: One per line
- **Comments**: Lines starting with `#` are ignored
- **Export prefix**: Optional `export` keyword before KEY=VALUE
- **Quote stripping**: Surrounding single or double quotes are removed from values
- **No variable expansion**: Variables like `$HOME` are treated as literal strings
- **No escape sequences**: Backslashes are treated as literal characters
- **No multi-line values**: Each line is treated as a single variable

## How It Works

`envee` reads the specified env file, parses the KEY=VALUE pairs, merges them with your current environment variables (env file variables override existing ones), and executes the child process with the merged environment.

The child process inherits the full environment (existing + injected variables) and `envee` passes through stdin/stdout/stderr, so it behaves transparently.
