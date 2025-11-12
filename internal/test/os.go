package test

import (
	"os/exec"
	"strings"
)

func RunFile(path string) ([]string, bool) {
	var cmd *exec.Cmd
	if len(path) > 0 {
		cmd = exec.Command("./wet", path)
	} else {
		cmd = exec.Command("./wet")
	}

	output, err := cmd.Output()
	if len(output) == 0 {
		return []string{}, err == nil
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), err == nil
}

