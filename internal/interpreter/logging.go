package interpreter

import (
	"fmt"
	"strings"
)

func (ip *Interpreter) runtimev(format string, args...any) bool {
	if !IsVerboseRuntime() {
		return false
	}

	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(colorYellow)
	sb.WriteString("Runtime: ")
	sb.WriteString(colorReset)
	sb.WriteString(format)
	fmt.Printf(sb.String(), args...)

	return false
}

func (ip *Interpreter) runtimeverr(format string, args...any) *StepResult {
	return &StepResult{
		status: ip.runtimev(format, args...),
		result: nil,
	}
}
