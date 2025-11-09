package tools

import (
	"fmt"
	"os"
)

func ToolLsf(dir string) (int64, error) {
	path, err := fixPath(dir)
	if err != nil {
		return 0, fmt.Errorf("failed to list directory files in %s: %w", dir, err)
	}

	list, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("failed to list directory files in %s: %w", dir, err)
	}

	var result int64

	for _, entry := range list {
		if entry.Type().IsRegular() {
			result++
		}
	}

	return result, nil
}
