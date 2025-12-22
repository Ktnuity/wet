package test

import (
	"os/exec"
	"runtime"
	"strings"
)

func RunFile(path string, forward []string) ([]string, bool) {
	args := make([]string, 0, len(forward) + 1)
	if len(path) > 0 {
		args = append(args, path)
	}
	args = append(args, forward...)
	executable := "./bin/wet"
	if runtime.GOOS == "windows" {
		executable = ".\\bin\\wet.exe"
	}
	cmd := exec.Command(executable, args...)
	output, err := cmd.Output()
	if len(output) == 0 {
		return []string{}, err == nil
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), err == nil
}

