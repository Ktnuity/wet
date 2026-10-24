package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RejectRootType int
const (
	RejectRoot RejectRootType = iota
	AllowRoot
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

func fixPath(path string, rootPolicy RejectRootType) (string, error) {
	if strings.HasPrefix(path, "/") {
		git, err := locateGit()
		if err != nil {
			return "", fmt.Errorf("failed to fix path '%s'. failed to locate git: %w", path, err)
		}

		git = filepath.Clean(git)
		result := filepath.Clean(git + path)
		if !strings.HasPrefix(result, git + string(filepath.Separator)) && (rootPolicy == RejectRoot || git != result) {
			if git == result {
				return "", fmt.Errorf("failed to fix path '%s'. cannot use root itself.", path)
			}
			return "", fmt.Errorf("failed to fix path '%s'. result path '%s' is outside git root '%s'.", path, result, git)
		}

		return result, nil
	} else if strings.HasPrefix(path, "./") {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}

		cwd = filepath.Clean(cwd)
		result := filepath.Clean(cwd + path[1:])
		if !strings.HasPrefix(result, cwd + string(filepath.Separator)) && (rootPolicy == RejectRoot || cwd != result) {
			if cwd == result {
				return "", fmt.Errorf("failed to fix path '%s'. cannot use root itself.", path)
			}
			return "", fmt.Errorf("failed to fix path '%s'. result path '%s' is outside cwd root '%s'.", path, result, cwd)
		}

		return result, nil
	} else if strings.HasPrefix(path, ":") {
		wetDir, err := getTokenDir()
		if err != nil {
			return "", fmt.Errorf("failed to fix token '%s': %w", path, err)
		}

		preDir := filepath.Clean(wetDir)
		result := filepath.Clean(preDir + "/" + path[1:])

		if !strings.HasPrefix(result, preDir + string(filepath.Separator)) && (rootPolicy == RejectRoot || preDir != result) {
			if preDir == result {
				return "", fmt.Errorf("failed to fix path '%s'. cannot use root itself.", path)
			}
			return "", fmt.Errorf("failed to fix path '%s'. result path '%s' is outside wet dir '%s'.", path, result, preDir)
		}

		parent := filepath.Dir(result)
		if err := os.MkdirAll(parent, 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}

		return result, nil
	} else {
		return "", fmt.Errorf("failed to fix path '%s'. invalid path type.", path)
	}
}
