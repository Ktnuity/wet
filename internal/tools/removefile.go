package tools

import (
	"fmt"
	"os"
)

func ToolRemoveFile(res string) error {
	path, err := fixPath(res, RejectRoot)
	if err != nil {
		return fmt.Errorf("failed to remove file %s: %w", res, err)
	}

	err = os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove file %s: %w", res, err)
	}

	return nil
}
