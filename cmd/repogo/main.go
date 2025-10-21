// RepoGo is a CLI tool that summarizes repository information.
// It can be used to generate summaries as input for LLM models, documentation, or analysis.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AndersonTsaiTW/RepoGo/internal/analyzer"
	"github.com/AndersonTsaiTW/RepoGo/internal/config"
	"github.com/AndersonTsaiTW/RepoGo/internal/git"
	"github.com/AndersonTsaiTW/RepoGo/internal/models"
	"github.com/AndersonTsaiTW/RepoGo/internal/renderer"
	"github.com/AndersonTsaiTW/RepoGo/internal/scanner"
)

func usage() {
	fmt.Fprintf(os.Stderr, `repogo %s

Usage:
  repogo [paths...] [flags]

Examples:
  repogo .
  repogo src main.go
  repogo . -o context.md
  repogo . --include "*.go,*.md" --exclude "*_test.go,vendor"

Flags:
`, config.Version)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	cfg := config.ParseFlags()

	if *cfg.Help {
		usage()
		return
	}
	if *cfg.Version {
		fmt.Println("repogo", config.Version)
		return
	}

	paths := config.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	rootAbs, _ := scanner.ResolveRoot(paths)

	doc := models.OutputDoc{Location: rootAbs}
	if gi, err := git.GetInfo(rootAbs); err == nil {
		doc.Git = gi
	} // Otherwise leave empty, will output "Not a git repository" later

	includes := scanner.SplitList(*cfg.Include)
	excludes := scanner.SplitList(*cfg.Exclude)

	files, structure := scanner.CollectFiles(rootAbs, paths, includes, excludes)
	doc.Structure = structure

	var totalTokens, totalLines, binaryCount, skippedByToken int

	for _, p := range files {
		rel, _ := filepath.Rel(rootAbs, p)
		if rel == "." {
			continue
		}
		info, err := os.Stat(p)
		if err != nil {
			doc.Files = append(doc.Files, models.FileEntry{
				Path:             filepath.ToSlash(rel),
				ReadErrorMessage: err.Error(),
			})
			continue
		}
		if info.IsDir() {
			continue
		}
		entry := models.FileEntry{
			Path: filepath.ToSlash(rel),
			Size: info.Size(),
		}

		f, err := os.Open(p)
		if err != nil {
			entry.ReadErrorMessage = err.Error()
			doc.Files = append(doc.Files, entry)
			continue
		}
		content, isBinary, truncated, lines := analyzer.ReadFileContent(f, *cfg.MaxFileSize)
		_ = f.Close()
		entry.IsBinary = isBinary
		entry.Truncated = truncated

		if !isBinary {
			entry.Content = string(content)
			entry.LanguageHint = analyzer.GuessLanguage(rel)
			totalLines += lines
		} else {
			binaryCount++
		}

		addTokens := analyzer.EstimateTokens(entry.Content)
		if *cfg.MaxTokens > 0 && totalTokens+addTokens > *cfg.MaxTokens {
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

	doc.Summary = models.Summary{
		TotalFiles:       len(doc.Files),
		TotalLines:       totalLines,
		EstimatedTokens:  totalTokens,
		SkippedByLimit:   skippedByToken,
		BinaryFilesCount: binaryCount,
	}

	var out bytes.Buffer
	switch strings.ToLower(*cfg.Format) {
	case "json":
		_ = renderer.RenderJSON(&out, doc)
	default:
		renderer.RenderMarkdown(&out, doc)
	}
	if *cfg.ShowTokens {
		fmt.Fprintf(&out, "\nEstimated tokens: %d\n", doc.Summary.EstimatedTokens)
	}

	if *cfg.Output == "" {
		_, _ = os.Stdout.Write(out.Bytes())
	} else {
		if err := os.WriteFile(*cfg.Output, out.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write output: %v\n", err)
			os.Exit(1)
		}
	}
}
