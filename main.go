package main

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"syscall"
)

var version = "dev"

func main() {
	var files stringSlice
	var overrides kvSlice
	var showVersion, dryRun, printEnv, validate, shell bool

	flag.Var(&files, "f", "path to env file (repeatable, default .env)")
	flag.Var(&overrides, "e", "set env var as KEY=VALUE (repeatable, highest precedence)")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&dryRun, "dry-run", false, "print resolved environment and exit")
	flag.BoolVar(&printEnv, "print", false, "alias for --dry-run")
	flag.BoolVar(&validate, "validate", false, "check env file(s) for issues and exit")
	flag.BoolVar(&shell, "shell", false, "spawn $SHELL with loaded environment")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f .env] [-e KEY=VALUE] [--dry-run|--validate|--shell] -- command [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("envee %s\n", version)
		os.Exit(0)
	}

	dryRun = dryRun || printEnv

	if len(files) == 0 {
		files = stringSlice{".env"}
	}

	modeCount := 0
	if dryRun {
		modeCount++
	}
	if validate {
		modeCount++
	}
	if shell {
		modeCount++
	}
	if modeCount > 1 {
		fmt.Fprintln(os.Stderr, "error: --dry-run, --validate, and --shell are mutually exclusive")
		os.Exit(2)
	}

	args := flag.Args()

	if modeCount == 0 && len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing command to run")
		flag.Usage()
		os.Exit(2)
	}

	if shell && len(args) > 0 {
		fmt.Fprintln(os.Stderr, "warning: extra arguments after -- ignored in --shell mode")
	}

	merged := make(map[string]string)
	var allWarnings []ParseWarning

	for _, file := range files {
		envFromFile, warnings, err := parseEnvFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read env file %q: %v\n", file, err)
			os.Exit(1)
		}
		allWarnings = append(allWarnings, warnings...)
		maps.Copy(merged, envFromFile)
	}

	if validate {
		if len(allWarnings) == 0 {
			fmt.Fprintln(os.Stderr, "no issues found")
			os.Exit(0)
		}
		for _, w := range allWarnings {
			fmt.Fprintln(os.Stderr, w)
		}
		os.Exit(1)
	}

	maps.Copy(merged, parseKVOverrides(overrides))

	mergedEnv := mergeEnv(os.Environ(), merged)

	if dryRun {
		for _, line := range envToSortedLines(mergedEnv) {
			fmt.Println(line)
		}
		os.Exit(0)
	}

	if shell {
		sh := os.Getenv("SHELL")
		if sh == "" {
			sh = "/bin/sh"
		}
		cmd := exec.Command(sh)
		cmd.Env = mergedEnv
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
			fmt.Fprintf(os.Stderr, "failed to run shell: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = mergedEnv
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		os.Exit(1)
	}
}
