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
	pathSrc, err := fixPath(src, AllowRoot)
	if err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	pathDst, err := fixPath(dst, AllowRoot)
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

	if err := copyRegularFile(pathSrc, pathDst, srcInfo.Mode()); err != nil {
		return false, fmt.Errorf("failed to copy file %s to %s: %w", src, dst, err)
	}

	return true, nil
}

func copyDirectory(pathSrc, pathDst, src, dst string) (bool, error) {
	entries, err := os.ReadDir(pathSrc)
	if err != nil {
		return false, fmt.Errorf("failed to copy directory %s to %s: %w", src, dst, err)
	}

	info, err := os.Stat(pathSrc)
	if err != nil {
		return false, fmt.Errorf("failed to copy directory %s to %s: %w", src, dst, err)
	}

	dstExists := false
	if _, err := os.Stat(pathDst); err == nil {
		dstExists = true
	} else {
		if err := os.MkdirAll(pathDst, 0700); err != nil {
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

			info, err := entry.Info()
			if err != nil {
				return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
			}

			if err := copyRegularFile(srcPath, dstPath, info.Mode()); err != nil {
				return dstExists, fmt.Errorf("failed to copy file %s to %s: %w", srcPath, dstPath, err)
			}
		}
	}

	if err := os.Chmod(pathDst, info.Mode().Perm()); err != nil {
		return dstExists, fmt.Errorf("failed to copy directory %s to %s: %w", src, dst, err)
	}

	return true, nil
}

func copyRegularFile(pathSrc, pathDst string, mode os.FileMode) error {
	srcFile, err := os.Open(pathSrc)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(pathDst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		return err
	}

	if err := dstFile.Close(); err != nil {
		return err
	}

	// OpenFile's mode is reduced by the umask, so set it explicitly
	return os.Chmod(pathDst, mode.Perm())
}
