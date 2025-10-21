// Package analyzer provides file content analysis functionality.
package analyzer

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
)

// ReadFileContent reads and analyzes a file, detecting if it's binary,
// checking for truncation, and counting lines.
func ReadFileContent(f *os.File, maxSize int) ([]byte, bool, bool, int) {
	defer f.Seek(0, 0) // Conservative approach: reset offset to zero after reading (though not used later)
	// Read max bytes to detect binary and do truncation
	buf := make([]byte, maxSize)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, false, false, 0
	}
	data := buf[:n]

	isBin := bytes.IndexByte(data, 0x00) >= 0 // Simple: treat files containing NUL as binary
	trunc := false
	if maxSize > 0 && int64(n) > int64(maxSize) {
		data = data[:maxSize]
		trunc = true
	}
	lines := bytes.Count(data, []byte{'\n'})
	return data, isBin, trunc, lines
}

// GuessLanguage returns the language identifier for syntax highlighting
// based on the file extension.
func GuessLanguage(filePath string) string {
	switch strings.ToLower(path.Ext(filePath)) {
	case ".go":
		return "go"
	case ".js", ".mjs", ".cjs", ".jsx":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".java":
		return "java"
	case ".cs":
		return "csharp"
	case ".c", ".h":
		return "c"
	case ".cpp", ".cc", ".cxx", ".hpp", ".hh":
		return "cpp"
	case ".sh", ".bash", ".zsh":
		return "bash"
	case ".yml", ".yaml":
		return "yaml"
	case ".sql":
		return "sql"
	case ".html", ".htm":
		return "html"
	case ".css", ".scss":
		return "css"
	default:
		return ""
	}
}
