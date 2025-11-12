package stdlib

//go:generate rm -rf ./std
//go:generate cp -r ../../wetstd ./std

import (
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

//go:embed std/*
var stdFS embed.FS

func GetContent() ([]*types.SourceSnippet, error) {
	files, err := listFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to load std: %v", err)
	}

	snippets := make([]*types.SourceSnippet, 0, 8)

	for _, file := range files {
		content, success, err := getFile(file)
		if success {
			lines := strings.Split(content, "\n")
			snippet := &types.SourceSnippet{
				Name: file,
				Start: 1,
				End: len(lines),
				Lines: make([]*types.SourceLine, 0, len(lines)),
			}

			/* how do I get the line number here? */
			for idx, line := range lines {
				snippet.Lines = append(snippet.Lines, &types.SourceLine{
					Parent: snippet,
					Content: line,
					Line: idx + 1,
				})
			}

			snippets = append(snippets, snippet)
		} else if err != nil {
			return nil, fmt.Errorf("failed to load std file '%s': %v", file, err)
		}
	}

	return snippets, nil
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

	return strings.TrimSuffix(string(data), "\n"), true, nil
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

	for line := range strings.Lines(string(data)) {
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
