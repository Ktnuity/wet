package test

import (
	"os"
	"strings"
)

func LoadTest() []string {
	data, err := os.ReadFile("./test.log")
	if err != nil {
		return []string{}
	}

	return strings.Split(string(data), "\n")
}

func GetTests() []string {
	data, err := os.ReadFile("./test.txt")
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
