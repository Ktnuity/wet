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

	srcInfo, err := os.Stat(pathSrc)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	if srcInfo.IsDir() {
		return copyDirectory(pathSrc, pathDst, src, dst)
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

func copyDirectory(pathSrc, pathDst, src, dst string) (bool, error) {
	entries, err := os.ReadDir(pathSrc)
	if err != nil {
		return false, fmt.Errorf("failed to copy directory %s to %s: %w", src, dst, err)
	}

	dstExists := false
	if _, err := os.Stat(pathDst); err == nil {
		dstExists = true
	} else {
		if err := os.MkdirAll(pathDst, 0755); err != nil {
			return false, fmt.Errorf("failed to copy directory %s to %s: %w", src, dst, err)
		}
	}

	for _, entry := range entries {
		srcPath := pathSrc + string(os.PathSeparator) + entry.Name()
		dstPath := pathDst + string(os.PathSeparator) + entry.Name()

		if entry.IsDir() {
			_, err := copyDirectory(srcPath, dstPath, src+"/"+entry.Name(), dst+"/"+entry.Name())
			if err != nil {
				return dstExists, err
			}
		} else {
			// For files in directories, override if exists
			if _, err := os.Stat(dstPath); err == nil {
				if err := os.Remove(dstPath); err != nil {
					return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
				}
			}

			srcFile, err := os.Open(srcPath)
			if err != nil {
				return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
			}

			dstFile, err := os.Create(dstPath)
			if err != nil {
				srcFile.Close()
				return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
			}

			_, err = io.Copy(dstFile, srcFile)
			srcFile.Close()
			dstFile.Close()
			if err != nil {
				return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
			}
		}
	}

	return true, nil
}
