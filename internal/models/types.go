// Package models defines the core data structures used throughout the application.
package models

// GitInfo contains information about the Git repository.
type GitInfo struct {
	Commit string `json:"commit"`
	Branch string `json:"branch"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

// FileEntry represents a single file in the repository.
type FileEntry struct {
	Path             string `json:"path"`
	Size             int64  `json:"size"`
	IsBinary         bool   `json:"is_binary"`
	Truncated        bool   `json:"truncated"`
	LanguageHint     string `json:"language_hint,omitempty"`
	Content          string `json:"content,omitempty"`
	ReadErrorMessage string `json:"read_error_message,omitempty"`
}

// Summary contains statistics about the scanned repository.
type Summary struct {
	TotalFiles       int `json:"total_files"`
	TotalLines       int `json:"total_lines"`
	EstimatedTokens  int `json:"estimated_tokens"`
	SkippedByLimit   int `json:"skipped_by_token_limit"`
	BinaryFilesCount int `json:"binary_files"`
}

// OutputDoc is the complete document structure for output.
type OutputDoc struct {
	Location  string      `json:"location"`
	Git       *GitInfo    `json:"git,omitempty"`
	Structure string      `json:"structure"`
	Files     []FileEntry `json:"files"`
	Summary   Summary     `json:"summary"`
}
