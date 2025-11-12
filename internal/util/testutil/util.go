package testutil

import (
	"fmt"
	"strings"
)

func Sprintf(format string, a...any) string {
	result := fmt.Sprintf(format, a...)

	mappings := map[string]string{
		"<red>": "\033[31m",
		"<green>": "\033[32m",
		"<gold>": "\033[33m",
		"<blue>": "\033[34m",
		"<purple>": "\033[35m",
		"<cyan>": "\033[36m",
		"<gray>": "\033[37m",
		"<lightgray>": "\033[90m",
		"<pink>": "\033[91m",
		"<lime>": "\033[92m",
		"<yellow>": "\033[93m",
		"<reset>": "\033[0m",
	}

	for key, value := range mappings {
		result = strings.ReplaceAll(result, key, value)
	}

	return result
}

func Printf(format string, a...any) {
	fmt.Print(Sprintf(format, a...))
}

func FindDiscrepency(match []string, current []string) []string {
	minLen := min(len(current), len(match))

	// (1) Check for discrepancy within common range
	for i := range minLen {
		if match[i] != current[i] {
			result := []string{}

			// Add previous 4 lines
			start := max(0, i-4)
			for j := start; j < i; j++ {
				result = append(result, Sprintf("<lime>%s<reset>", match[j]))
			}

			// Add separator
			//result = append(result, "\033[37m----------------------------------------\033[0m")

			// Add differing lines
			result = append(result, Sprintf("<lightgray>+%s<reset>", match[i]))
			result = append(result, Sprintf("<pink>-%s<reset>", current[i]))

			return result
		}
	}

	// (2) If current runs out but match has more, return empty
	if len(current) < len(match) {
		return []string{}
	}

	// (3) If match runs out but current has more
	if len(match) < len(current) {
		result := []string{}

		// Add previous 4 lines from match
		start := max(0, len(match)-4)
		for j := start; j < len(match); j++ {
			result = append(result, Sprintf("<lime>%s<reset>", match[j]))
		}

		// Add separator
		//result = append(result, "\033[37m----------------------------------------\033[0m")

		// Add last line of match
		result = append(result, Sprintf("<pink>- %s<reset>", match[len(match)-1]))

		// Add up to 4 lines after match from current
		end := min(len(current), len(match)+4)
		for j := len(match); j < end; j++ {
			result = append(result, Sprintf("<lightgray>%s<reset>", current[j]))
		}

		return result
	}

	return []string{}
}
