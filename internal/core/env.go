// internal/core/env.go
package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	const maxLineSize = 1024 * 1024
	buf := make([]byte, maxLineSize)

	scanner := bufio.NewScanner(f)
	scanner.Buffer(buf, maxLineSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// handle export prefix
		line = strings.TrimPrefix(line, "export ")

		key, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("line %d: missing '='", lineNum)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// validate key
		if key == "" || strings.ContainsAny(key, " \t") {
			return fmt.Errorf("line %d: invalid key %q", lineNum, key)
		}

		// strip inline comments (only outside quotes)
		if !isQuoted(value) {
			if idx := strings.Index(value, " #"); idx != -1 {
				value = strings.TrimSpace(value[:idx])
			}
		}

		// strip surrounding quotes
		value = stripQuotes(value)

		// real environment always wins
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func isQuoted(s string) bool {
	return len(s) >= 2 &&
		((s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\''))
}

func stripQuotes(s string) string {
	if isQuoted(s) {
		return s[1 : len(s)-1]
	}
	return s
}
