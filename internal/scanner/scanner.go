// Package scanner provides file system scanning and filtering functionality.
package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// CollectFiles scans the filesystem and collects files based on include/exclude patterns.
// Returns a list of file paths and a string representation of the directory structure.
func CollectFiles(root string, inputs, includes, excludes []string) ([]string, string) {
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

// ResolveRoot determines the root directory from the given input paths.
// Uses first directory as root; if all are files, finds common parent.
func ResolveRoot(inputs []string) (string, error) {
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
