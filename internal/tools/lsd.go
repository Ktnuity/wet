package tools

import (
	"fmt"
	"os"
)

func ToolLsd(dir string) (int64, error) {
	path, err := fixPath(dir)
	if err != nil {
		return 0, fmt.Errorf("failed to list directory sub-dirs in %s: %w", dir, err)
	}

	list, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("failed to list directory sub-dirs in %s: %w", dir, err)
	}

	var result int64

	for _, entry := range list {
		if entry.Type().IsDir() {
			result++
		}
	}

	return result, nil
}
