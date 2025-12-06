package main

import (
	"bufio"
	"os"
	"strings"
)

// parseEnvFile reads simple KEY=VALUE lines, ignoring empty lines and lines starting with '#'
// Surrounding single or double quotes on VALUE are stripped if present.
// No expansion, no escapes, no multi-line support—kept intentionally simple.
func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
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
		if idx <= 0 {
			// Skip invalid line quietly; alternatively, warn:
			// fmt.Fprintf(os.Stderr, "warning: ignoring invalid line %d: %q\n", lineNo, line)
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip matching surrounding quotes, if any.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if key != "" {
			result[key] = val
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return result, nil
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
	for k, v := range overlay {
		m[k] = v
	}
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}
