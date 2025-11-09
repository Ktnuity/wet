package tools

import (
	"fmt"
	"io"
	"os"
)

// Returns (true, ###) if dst path is occupied once the function returns.
// Returns (###, error) if the copy failed.
// Returning (true, error) means the copy failed, because a file already exist in the destination path.
// Returning (false, error) means the copy failed, but not because the file already exist.
// Returning (true, nil) means the copy succeeded and a file now exist in the destination path.
// Returning (false, nil) makes no sense.
func ToolCopyFile(src, dst string) (bool, error) {
	pathSrc, err := fixPath(src)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	pathDst, err := fixPath(dst)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	if _, err := os.Stat(pathDst); err == nil {
		return true, fmt.Errorf("failed to copy file %s: %s already exist", src, dst)
	}

	srcFile, err := os.Open(pathSrc)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(pathDst)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	return true, nil
}
