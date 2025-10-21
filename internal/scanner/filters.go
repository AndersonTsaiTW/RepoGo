// Package scanner provides file system scanning and filtering functionality.
package scanner

import "strings"

// SplitList splits a comma-separated string into a slice of trimmed strings.
// Returns nil if the input is empty or contains only whitespace.
func SplitList(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
