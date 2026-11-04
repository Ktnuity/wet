package main

import (
	"io"
	"os"
	"path/filepath"
)

func main() {
	cdToGitRoot()

	dest := filepath.Join("internal", "stdlib", "std")
	src := "wetstd"

	os.RemoveAll(dest)

	if err := copyDir(src, dest); err != nil {
		panic(err)
	}
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func cdToGitRoot() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			if err := os.Chdir(cwd); err != nil {
				panic(err)
			}
			break
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			panic("no .git directory found")
		}
		cwd = parent
	}
}
