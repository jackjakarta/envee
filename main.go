package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var file string

	flag.StringVar(&file, "f", ".env", "path to env file (default .env)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f .env] -- command [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing command to run")
		flag.Usage()
		os.Exit(2)
	}

	envFromFile, err := parseEnvFile(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read env file %q: %v\n", file, err)
		os.Exit(1)
	}

	mergedEnv := mergeEnv(os.Environ(), envFromFile)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = mergedEnv
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// If the process ran but exited non-zero, propagate that exit code.
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}

		// Otherwise it likely failed to start.
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		os.Exit(1)
	}
}
