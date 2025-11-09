package tools

import (
	"fmt"
	"os"
)

func ToolMakeDirectory(res string) error {
	path, err := fixPath(res)
	if err != nil {
		return fmt.Errorf("failed to make directory %s: %w", res, err)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to make directory %s: %w", res, err)
	}

	return nil
}
