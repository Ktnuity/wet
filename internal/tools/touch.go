package tools

import (
	"fmt"
	"os"
	"time"
)

func ToolTouchFile(res string) error {
	path, err := fixPath(res)
	if err != nil {
		return fmt.Errorf("failed to touch file %s: %w", res, err)
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to touch file %s: %w", res, err)
	}
	file.Close()

	now := time.Now()
	err = os.Chtimes(path, now, now)
	if err != nil {
		return fmt.Errorf("failed to touch file %s: %w", res, err)
	}

	return nil
}
