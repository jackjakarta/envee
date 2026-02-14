package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"sort"
	"strings"
)

// stringSlice implements flag.Value for repeatable -f flags.
type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ", ") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// kvSlice implements flag.Value for repeatable -e KEY=VALUE flags.
type kvSlice []string

func (k *kvSlice) String() string { return strings.Join(*k, ", ") }
func (k *kvSlice) Set(v string) error {
	if !strings.Contains(v, "=") {
		return fmt.Errorf("invalid format %q: expected KEY=VALUE", v)
	}
	*k = append(*k, v)
	return nil
}

// ParseWarning describes a non-fatal issue found while parsing an env file.
type ParseWarning struct {
	File   string
	Line   int
	Text   string
	Reason string
}

func (w ParseWarning) String() string {
	return fmt.Sprintf("%s:%d: %s (line: %s)", w.File, w.Line, w.Reason, w.Text)
}

// parseEnvFile reads simple KEY=VALUE lines, ignoring empty lines and lines starting with '#'.
// Surrounding single or double quotes on VALUE are stripped if present.
// No expansion, no escapes, no multi-line support—kept intentionally simple.
func parseEnvFile(path string) (map[string]string, []ParseWarning, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, nil, err
	}

	defer f.Close()

	result := make(map[string]string)
	var warnings []ParseWarning
	sc := bufio.NewScanner(f)
	lineNo := 0

	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Optional: allow leading "export "
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		// Split on first '='
		idx := strings.IndexByte(line, '=')

		if idx < 0 {
			warnings = append(warnings, ParseWarning{
				File: path, Line: lineNo, Text: sc.Text(), Reason: "missing '=' delimiter",
			})
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if key == "" {
			warnings = append(warnings, ParseWarning{
				File: path, Line: lineNo, Text: sc.Text(), Reason: "empty key",
			})
			continue
		}

		// Strip matching surrounding quotes, if any.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if _, exists := result[key]; exists {
			warnings = append(warnings, ParseWarning{
				File: path, Line: lineNo, Text: sc.Text(), Reason: fmt.Sprintf("duplicate key %q", key),
			})
		}
		result[key] = val
	}

	if err := sc.Err(); err != nil {
		return nil, nil, err
	}

	return result, warnings, nil
}

// parseKVOverrides converts a slice of KEY=VALUE strings into a map.
// Last value wins for duplicate keys.
func parseKVOverrides(kvs []string) map[string]string {
	m := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		if idx := strings.IndexByte(kv, '='); idx > 0 {
			m[kv[:idx]] = kv[idx+1:]
		}
	}
	return m
}

// envToSortedLines sorts a KEY=VALUE slice alphabetically.
func envToSortedLines(env []string) []string {
	sorted := make([]string, len(env))
	copy(sorted, env)
	sort.Strings(sorted)
	return sorted
}

// mergeEnv overlays the key/values from overlay onto base (os.Environ()),
// returning a new []string in KEY=VALUE form.
func mergeEnv(base []string, overlay map[string]string) []string {
	m := make(map[string]string, len(base)+len(overlay))

	for _, kv := range base {
		if i := strings.IndexByte(kv, '='); i > 0 {
			m[kv[:i]] = kv[i+1:]
		}
	}

	maps.Copy(m, overlay)
	out := make([]string, 0, len(m))

	for k, v := range m {
		out = append(out, k+"="+v)
	}

	return out
}
