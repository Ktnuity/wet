package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ToolUnzipResult struct {
	DirCount		int64
	FileCount		int64
}

func ToolUnzipFile(dst, res string) (ToolUnzipResult, error) {
	var result ToolUnzipResult

	pathRes, err := fixPath(res)
	if err != nil {
		return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
	}

	pathDst, err := fixPath(dst)
	if err != nil {
		return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
	}

	reader, err := zip.OpenReader(pathRes)
	if err != nil {
		return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
	}
	defer reader.Close()

	if err := os.MkdirAll(pathDst, 0755); err != nil {
		return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
	}

	seenDirs := make(map[string]bool)
	seenFiles := make(map[string]bool)

	for _, file := range reader.File {
		path := filepath.Join(pathDst, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(pathDst)+string(os.PathSeparator)) {
			return result, fmt.Errorf("failed to unzip file %s: illegal file path: %s", res, file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
			}

			rootDir := strings.Split(strings.TrimPrefix(file.Name, "/"), "/")[0]
			if rootDir != "" && !seenDirs[rootDir] {
				seenDirs[rootDir] = true
				result.DirCount++
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
		}

		fileReader, err := file.Open()
		if err != nil {
			dstFile.Close()
			return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
		}

		_, err = io.Copy(dstFile, fileReader)
		fileReader.Close()
		dstFile.Close()

		if err != nil {
			return result, fmt.Errorf("failed to unzip file %s: %w", res, err)
		}

		rootItem := strings.Split(strings.TrimPrefix(file.Name, "/"), "/")[0]
		if rootItem != "" && !seenFiles[rootItem] {
			seenFiles[rootItem] = true
			result.FileCount++
		}
	}

	return result, nil
}
