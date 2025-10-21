// Package analyzer provides file content analysis functionality.
package analyzer

// EstimateTokens provides a rough estimate of the number of tokens in a string.
// Uses a simple heuristic of approximately 4 characters per token.
func EstimateTokens(s string) int {
	// Very rough estimate: ~4 chars/token
	return (len(s) + 3) / 4
}
