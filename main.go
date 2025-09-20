package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	Version       = "0.1.0-pure"
	flagVersion   = flag.Bool("v", false, "print version")
	flagHelp      = flag.Bool("h", false, "show help")
	flagOutput    = flag.String("o", "", "output file (default stdout)")
	flagInclude   = flag.String("include", "", "comma-separated glob(s) to include (supports *, ?, [class])")
	flagExclude   = flag.String("exclude", "", "comma-separated glob(s) to exclude (supports *, ?, [class])")
	flagFormat    = flag.String("format", "markdown", "output format: markdown|json")
	flagTokens    = flag.Bool("tokens", false, "print estimated token count")
	flagMaxFileSz = flag.Int("max-file-size", 16*1024, "per-file size limit in bytes before truncation")
	flagMaxTokens = flag.Int("max-tokens", 0, "stop when total estimated tokens reach this number (0 = no limit)")
)

type GitInfo struct {
	Commit string `json:"commit"`
	Branch string `json:"branch"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

type FileEntry struct {
	Path             string `json:"path"`
	Size             int64  `json:"size"`
	IsBinary         bool   `json:"is_binary"`
	Truncated        bool   `json:"truncated"`
	LanguageHint     string `json:"language_hint,omitempty"`
	Content          string `json:"content,omitempty"`
	ReadErrorMessage string `json:"read_error_message,omitempty"`
}

type Summary struct {
	TotalFiles       int `json:"total_files"`
	TotalLines       int `json:"total_lines"`
	EstimatedTokens  int `json:"estimated_tokens"`
	SkippedByLimit   int `json:"skipped_by_token_limit"`
	BinaryFilesCount int `json:"binary_files"`
}

type OutputDoc struct {
	Location  string      `json:"location"`
	Git       *GitInfo    `json:"git,omitempty"`
	Structure string      `json:"structure"`
	Files     []FileEntry `json:"files"`
	Summary   Summary     `json:"summary"`
}

func usage() {
	fmt.Fprintf(os.Stderr, `repopack %s (pure stdlib)

Usage:
  repopack [paths...] [flags]

Examples:
  repopack .
  repopack src main.go
  repopack . -o context.md
  repopack . --include "*.go,*.md" --exclude "*_test.go,vendor"

Flags:
`, Version)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *flagHelp {
		usage()
		return
	}
	if *flagVersion {
		fmt.Println("repopack", Version)
		return
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	rootAbs, _ := resolveRoot(paths)

	doc := OutputDoc{Location: rootAbs}
	if gi, err := getGitInfoViaCLI(rootAbs); err == nil {
		doc.Git = gi
	} // Otherwise leave empty, will output "Not a git repository" later

	includes := splitList(*flagInclude)
	excludes := splitList(*flagExclude)

	files, structure := collectFiles(rootAbs, paths, includes, excludes)
	doc.Structure = structure

	var totalTokens, totalLines, binaryCount, skippedByToken int

	for _, p := range files {
		rel, _ := filepath.Rel(rootAbs, p)
		if rel == "." {
			continue
		}
		info, err := os.Stat(p)
		if err != nil {
			doc.Files = append(doc.Files, FileEntry{
				Path:             filepath.ToSlash(rel),
				ReadErrorMessage: err.Error(),
			})
			continue
		}
		if info.IsDir() {
			continue
		}
		entry := FileEntry{
			Path: filepath.ToSlash(rel),
			Size: info.Size(),
		}

		f, err := os.Open(p)
		if err != nil {
			entry.ReadErrorMessage = err.Error()
			doc.Files = append(doc.Files, entry)
			continue
		}
		content, isBinary, truncated, lines := readFileContent(f, *flagMaxFileSz)
		_ = f.Close()
		entry.IsBinary = isBinary
		entry.Truncated = truncated

		if !isBinary {
			entry.Content = string(content)
			entry.LanguageHint = guessFence(rel)
			totalLines += lines
		} else {
			binaryCount++
		}

		addTokens := estimateTokens(entry.Content)
		if *flagMaxTokens > 0 && totalTokens+addTokens > *flagMaxTokens {
			skippedByToken++
			entry.Content = ""
			entry.Truncated = true
			entry.ReadErrorMessage = fmt.Sprintf("omitted due to --max-tokens budget (would add ~%d tokens)", addTokens)
			doc.Files = append(doc.Files, entry)
			break
		}
		totalTokens += addTokens
		doc.Files = append(doc.Files, entry)
	}

	doc.Summary = Summary{
		TotalFiles:       len(doc.Files),
		TotalLines:       totalLines,
		EstimatedTokens:  totalTokens,
		SkippedByLimit:   skippedByToken,
		BinaryFilesCount: binaryCount,
	}

	var out bytes.Buffer
	switch strings.ToLower(*flagFormat) {
	case "json":
		enc := json.NewEncoder(&out)
		enc.SetIndent("", "  ")
		_ = enc.Encode(doc)
	default:
		renderMarkdown(&out, doc)
	}
	if *flagTokens {
		fmt.Fprintf(&out, "\nEstimated tokens: %d\n", doc.Summary.EstimatedTokens)
	}

	if *flagOutput == "" {
		_, _ = os.Stdout.Write(out.Bytes())
	} else {
		if err := os.WriteFile(*flagOutput, out.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write output: %v\n", err)
			os.Exit(1)
		}
	}
}

// ---------- helpers (pure stdlib) ----------

func resolveRoot(inputs []string) (string, error) {
	// Use first directory as root; if all are files, find common parent
	var dirs []string
	for _, in := range inputs {
		ap, _ := filepath.Abs(in)
		if fi, err := os.Stat(ap); err == nil && fi.IsDir() {
			return ap, nil
		}
		dirs = append(dirs, filepath.Dir(ap))
	}
	if len(dirs) == 0 {
		return filepath.Abs(".")
	}
	base := dirs[0]
	for _, d := range dirs[1:] {
		base = commonBase(base, d)
	}
	return base, nil
}

func commonBase(a, b string) string {
	pa := strings.Split(filepath.Clean(a), string(os.PathSeparator))
	pb := strings.Split(filepath.Clean(b), string(os.PathSeparator))
	i := 0
	for i < len(pa) && i < len(pb) && pa[i] == pb[i] {
		i++
	}
	return filepath.Join(pa[:i]...)
}

func splitList(s string) []string {
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

func matchAny(patterns []string, rel string) bool {
	if len(patterns) == 0 {
		return false
	}
	rel = filepath.ToSlash(rel)
	for _, pat := range patterns {
		ok, _ := path.Match(pat, rel)
		if ok {
			return true
		}
		// Also try matching filename only
		ok, _ = path.Match(pat, path.Base(rel))
		if ok {
			return true
		}
	}
	return false
}

func collectFiles(root string, inputs, includes, excludes []string) ([]string, string) {
	seen := map[string]struct{}{}
	var files []string
	add := func(p string) {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			files = append(files, p)
		}
	}

	shouldKeep := func(rel string) bool {
		if matchAny(excludes, rel) {
			return false
		}
		if len(includes) == 0 {
			return true
		}
		return matchAny(includes, rel)
	}

	for _, in := range inputs {
		ap, err := filepath.Abs(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", in, err)
			continue
		}
		info, err := os.Stat(ap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", in, err)
			continue
		}
		if info.IsDir() {
			filepath.WalkDir(ap, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Fprintf(os.Stderr, "walk error %s: %v\n", p, err)
					return nil
				}
				rel, _ := filepath.Rel(root, p)
				if rel == "." {
					return nil
				}
				if !shouldKeep(rel) {
					if d.IsDir() {
						return fs.SkipDir
					}
					return nil
				}
				add(p)
				return nil
			})
		} else {
			rel, _ := filepath.Rel(root, ap)
			if shouldKeep(rel) {
				add(ap)
			}
		}
	}
	sort.Strings(files)
	structure := buildTree(root, files)
	return files, structure
}

func buildTree(root string, files []string) string {
	type node struct {
		name     string
		children map[string]*node
		file     bool
	}
	rootNode := &node{children: map[string]*node{}}
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		if rel == "." {
			continue
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		cur := rootNode
		for i, part := range parts {
			if cur.children[part] == nil {
				cur.children[part] = &node{name: part, children: map[string]*node{}}
			}
			cur = cur.children[part]
			if i == len(parts)-1 {
				cur.file = true
			}
		}
	}
	var b strings.Builder
	var walk func(*node, int)
	walk = func(n *node, depth int) {
		keys := make([]string, 0, len(n.children))
		for k := range n.children {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ch := n.children[k]
			b.WriteString(strings.Repeat("  ", depth))
			if ch.file {
				b.WriteString(k + "\n")
			} else {
				b.WriteString(k + "/\n")
			}
			walk(ch, depth+1)
		}
	}
	walk(rootNode, 0)
	return "```\n" + b.String() + "```"
}

func readFileContent(f *os.File, max int) ([]byte, bool, bool, int) {
	defer f.Seek(0, 0) // Conservative approach: reset offset to zero after reading (though not used later)
	// Read max bytes to detect binary and do truncation
	buf := make([]byte, max)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, false, false, 0
	}
	data := buf[:n]

	isBin := bytes.IndexByte(data, 0x00) >= 0 // Simple: treat files containing NUL as binary
	trunc := false
	if max > 0 && int64(n) > int64(max) {
		data = data[:max]
		trunc = true
	}
	lines := bytes.Count(data, []byte{'\n'})
	return data, isBin, trunc, lines
}

func guessFence(rel string) string {
	switch strings.ToLower(path.Ext(rel)) {
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

func estimateTokens(s string) int {
	// Very rough estimate: ~4 chars/token
	return (len(s) + 3) / 4
}

func renderMarkdown(w io.Writer, doc OutputDoc) {
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

// ---- Git via system git (pure stdlib) ----

func getGitInfoViaCLI(root string) (*GitInfo, error) {
	// Commands will fail if root is not a git repo
	run := func(args ...string) (string, error) {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}

	// First check if .git or HEAD is readable
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		// Could be in subdirectory, try git rev-parse
		if _, err2 := run("rev-parse", "--git-dir"); err2 != nil {
			return nil, fmt.Errorf("not a git repo")
		}
	}

	commit, err := run("rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	branch, _ := run("rev-parse", "--abbrev-ref", "HEAD")
	authorName, _ := run("log", "-1", "--pretty=%an")
	authorEmail, _ := run("log", "-1", "--pretty=%ae")
	dateRaw, _ := run("log", "-1", "--pretty=%ad", "--date=rfc")
	if dateRaw == "" {
		dateRaw = time.Now().Format(time.RFC1123Z)
	}
	return &GitInfo{
		Commit: commit,
		Branch: branch,
		Author: fmt.Sprintf("%s <%s>", authorName, authorEmail),
		Date:   dateRaw,
	}, nil
}
