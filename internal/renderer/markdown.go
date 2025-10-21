// Package renderer provides output rendering functionality for different formats.
package renderer

import (
	"fmt"
	"io"

	"github.com/AndersonTsaiTW/RepoGo/internal/models"
)

// RenderMarkdown renders the output document in Markdown format.
func RenderMarkdown(w io.Writer, doc models.OutputDoc) {
	fmt.Fprintln(w, "# Repository Context\n")
	fmt.Fprintln(w, "## File System Location\n")
	fmt.Fprintln(w, doc.Location, "\n")

	fmt.Fprintln(w, "## Git Info\n")
	if doc.Git == nil {
		fmt.Fprintln(w, "- Not a git repository\n")
	} else {
		fmt.Fprintf(w, "- Commit: %s\n", doc.Git.Commit)
		fmt.Fprintf(w, "- Branch: %s\n", doc.Git.Branch)
		fmt.Fprintf(w, "- Author: %s\n", doc.Git.Author)
		fmt.Fprintf(w, "- Date: %s\n\n", doc.Git.Date)
	}

	fmt.Fprintln(w, "## Structure")
	fmt.Fprintln(w, doc.Structure, "\n")

	fmt.Fprintln(w, "## File Contents\n")
	for _, f := range doc.Files {
		fmt.Fprintf(w, "### File: %s\n", f.Path)
		if f.ReadErrorMessage != "" && f.Content == "" && !f.IsBinary {
			fmt.Fprintf(w, "_Error: %s_\n\n", f.ReadErrorMessage)
			continue
		}
		if f.IsBinary {
			fmt.Fprintf(w, "_Binary file (size: %d bytes) — metadata only._\n\n", f.Size)
			continue
		}
		fmt.Fprintf(w, "```%s\n%s\n```\n\n", f.LanguageHint, f.Content)
		if f.Truncated {
			fmt.Fprintln(w, "_[truncated]_\n")
		}
		if f.ReadErrorMessage != "" {
			fmt.Fprintf(w, "_Note: %s_\n\n", f.ReadErrorMessage)
		}
	}

	fmt.Fprintln(w, "## Summary")
	fmt.Fprintf(w, "- Total files: %d\n", doc.Summary.TotalFiles)
	fmt.Fprintf(w, "- Total lines: %d\n", doc.Summary.TotalLines)
	fmt.Fprintf(w, "- Estimated tokens: %d\n", doc.Summary.EstimatedTokens)
	if doc.Summary.SkippedByLimit > 0 {
		fmt.Fprintf(w, "- Skipped due to token limit: %d file(s)\n", doc.Summary.SkippedByLimit)
	}
	if doc.Summary.BinaryFilesCount > 0 {
		fmt.Fprintf(w, "- Binary files detected: %d\n", doc.Summary.BinaryFilesCount)
	}
}
