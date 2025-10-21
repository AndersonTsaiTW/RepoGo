// Package git provides functionality to retrieve Git repository information.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AndersonTsaiTW/RepoGo/internal/models"
)

// GetInfo retrieves Git repository information from the specified root directory.
// Returns nil if the directory is not a Git repository.
func GetInfo(root string) (*models.GitInfo, error) {
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
	return &models.GitInfo{
		Commit: commit,
		Branch: branch,
		Author: fmt.Sprintf("%s <%s>", authorName, authorEmail),
		Date:   dateRaw,
	}, nil
}
