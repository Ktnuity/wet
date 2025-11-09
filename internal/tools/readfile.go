package tools

import (
	"fmt"
	"os"
)

func ToolReadfile(src string) (string, error) {
	path, err := fixPath(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", src, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", src, err)
	}

	return string(data), nil
}
