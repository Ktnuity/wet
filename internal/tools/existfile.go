package tools

import (
	"fmt"
	"os"
)

func ToolExistFile(res string) error {
	path, err := fixPath(res)
	if err != nil {
		return fmt.Errorf("failed to exists file %s: %w", res, err)
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("failed to exists file %s: %w", res, err)
	}

	return nil
}
