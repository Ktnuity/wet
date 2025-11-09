package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ToolDownload(url, dst string) error {
	path, err := fixPath(dst)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", dst, err)
	}

	return nil
}
