package tools

import (
	"fmt"
	"os"
)

func ToolMoveFile(src, dst string) error {
	pathSrc, err := fixPath(src)
	if err != nil {
		return fmt.Errorf("failed to move file %s to %s: %w", src, dst, err)
	}

	pathDst, err := fixPath(dst)
	if err != nil {
		return fmt.Errorf("failed to move file %s to %s: %w", src, dst, err)
	}

	err = os.Rename(pathSrc, pathDst)
	if err != nil {
		return fmt.Errorf("failed to move file %s to %s: %w", src, dst, err)
	}

	return nil
}
