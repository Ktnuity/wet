package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func locateGit() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := cwd
	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .git directory found")
		}
		dir = parent
	}
}

func getTokenDir() (string, error) {
	git, err := locateGit()
	if err != nil {
		return "", fmt.Errorf("failed to get token dir. failed to locate git: %w", err)
	}

	wetDir := git + "/.wet"
	s, err := os.Stat(wetDir)
	if err == nil {
		if s.IsDir() {
			return wetDir, nil
		} else {
			return "", fmt.Errorf("failed to get token dir. ./.wet is not a directory.")
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	wetDir = filepath.Join(homeDir, ".wet")
	s, err = os.Stat(wetDir)
	if err != nil {
		if err := os.MkdirAll(wetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create ~/.wet directory: %w", err)
		}

		return wetDir, nil
	}

	if !s.IsDir() {
		return "", fmt.Errorf("failed to get token dir. ~/.wet is not a directory.")
	}

	return wetDir, nil
}

func fixPath(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		git, err := locateGit()
		if err != nil {
			return "", fmt.Errorf("failed to fix path '%s'. failed to locate git: %w", path, err)
		}

		return git + path, nil
	} else if strings.HasPrefix(path, "./") {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}

		return cwd + path[1:], nil
	} else if strings.HasPrefix(path, ":") {
		wetDir, err := getTokenDir()
		if err != nil {
			return "", fmt.Errorf("failed to fix token '%s': %w", path, err)
		}

		fullPath := wetDir + "/" + path
		parent := filepath.Dir(fullPath)

		if err := os.MkdirAll(parent, 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}

		return fullPath, nil
	} else {
		return "", fmt.Errorf("failed to fix path '%s'. invalid path type.", path)
	}
}
