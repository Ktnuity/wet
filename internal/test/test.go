package test

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func GetTestFile() (string, error) {
	if runtime.GOOS == "linux" {
		return "./tests/_config/linux.test.log", nil
	} else if runtime.GOOS == "darwin" {
		return "./tests/_config/darwin.test.log", nil
	} else if runtime.GOOS == "windows" {
		return "./tests/_config/win32.test.log", nil
	} else {
		return "", fmt.Errorf("unsupported os name '%s'. supports one of 'linux', 'darwin', 'windows' for runtime.GOOS.", runtime.GOOS)
	}
}

func LoadTest() ([]string, error) {
	path, err := GetTestFile()
	if err != nil {
		return []string{}, fmt.Errorf("failed to load test, failed to get test file path: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		return []string{}, nil // reading an existing file could fail, that's an error, but when file doesn't exist,
							   // then that means we should still run the test to generate a new one. thus this
							   // isn't an error, and nil error is returned.
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return []string{}, fmt.Errorf("failed to load test, loading test file failed: %v", err)
	}

	return strings.Split(string(data), "\n"), nil
}

func GetTests() []string {
	data, err := os.ReadFile("./tests/_config/test.txt")
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(data), "\n")
	result := []string{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			result = append(result, trimmed)
		}
	}

	result = append(result, "")

	return result
}
