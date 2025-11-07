package stdlib

//go:generate rm -rf ./std
//go:generate cp -r ../../wetstd ./std

import (
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

//go:embed std/*
var stdFS embed.FS

func GetContent() (string, error) {
	files, err := listFiles()
	if err != nil {
		return "", fmt.Errorf("failed to load std: %v", err)
	}

	parts := make([]string, 0, len(files))

	for _, file := range files {
		content, success, err := getFile(file)
		if success {
			parts = append(parts, content)
		} else if err != nil {
			return "", fmt.Errorf("failed to load std file '%s': %v", file, err)
		}
	}

	return strings.Join(parts, "\n"), nil
}

func getFile(fileName string) (string, bool, error) {
	if !validFilename(fileName) {
		return "", false, nil
	}

	stdPath := "std/" + fileName
	_, err := fs.Stat(stdFS, stdPath)
	if err != nil {
		return "", false, nil
	}

	data, err := fs.ReadFile(stdFS, stdPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to open std file '%s': %v", stdPath, err)
	}

	return string(data), true, nil
}

func listFiles() ([]string, error) {
	result := make([]string, 0, 10)

	_, err := fs.Stat(stdFS, "std/std.txt")
	if err != nil {
		return result, fmt.Errorf("failed to stat std lib index.")
	}

	data, err := fs.ReadFile(stdFS, "std/std.txt")
	if err != nil {
		return result, fmt.Errorf("failed to open std lib index: %v", err)
	}

	input := strings.Split(string(data), "\n")

	for _, line := range input {
		fileName := strings.TrimSpace(line)
		if len(fileName) == 0 {
			continue
		}

		if !validFilename(fileName) {
			continue
		}

		result = append(result, fileName)
	}

	return result, nil
}

func validFilename(fileName string) bool {
	matched, _ := regexp.MatchString(`^[a-z]+\.wet$`, fileName)
	return matched
}
