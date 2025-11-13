package test

import (
	"os/exec"
	"strings"
)

func RunFile(path string, forward []string) ([]string, bool) {
	args := make([]string, 0, len(forward) + 1)
	if len(path) > 0 {
		args = append(args, path)
	}
	args = append(args, forward...)
	cmd := exec.Command("./wet", args...)
	output, err := cmd.Output()
	if len(output) == 0 {
		return []string{}, err == nil
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), err == nil
}

