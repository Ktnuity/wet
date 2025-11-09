package tools

import (
	"fmt"
	"os"
)

func ToolGetf(idx int, dir string) (string, error) {
	path, err := fixPath(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get file name in %s: %w", dir, err)
	}

	list, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file name in %s: %w", dir, err)
	}

	if idx < 0 {
		return "", fmt.Errorf("failed to get file name in %s: index(%d) out of bounds", dir, idx)
	}

	var inc int
	for _, entry := range list {
		if entry.Type().IsRegular() {
			if inc == idx {
				return entry.Name(), nil
			} else {
				inc++
			}
		}
	}

	return "", fmt.Errorf("failed to get file name in %s: index(%d) out of bounds(%d)", dir, idx, inc)
}
