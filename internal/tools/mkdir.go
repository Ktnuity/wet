package tools

import (
	"fmt"
	"os"
)

func ToolMakeDirectory(res string) error {
	if err := unavailable("mkdir"); err != nil {
		return err
	}

	path, err := fixPath(res, AllowRoot)
	if err != nil {
		return fmt.Errorf("failed to make directory %s: %w", res, err)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to make directory %s: %w", res, err)
	}

	return nil
}
