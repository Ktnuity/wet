package interpreter

import (
	"fmt"
	"strings"

	"github.com/ktnuity/wet/internal/errors"
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

func (ip *Interpreter) basicRuntimeverr(format string, args...any) *StepResult {
	return &StepResult{
		status: ip.runtimev(format, args...),
		result: fmt.Errorf(format, args...),
	}
}

func (ip *Interpreter) runtimevwarn(format string, args...any) bool {
	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(colorYellow)
	sb.WriteString("Warn: ")
	sb.WriteString(colorReset)
	sb.WriteString(format)
	fmt.Printf(sb.String(), args...)

	return false
}

func (ip *Interpreter) runtimeverr(d *StepData, format string, args...any) *StepResult {
	if d == nil {
		return ip.basicRuntimeverr(format, args...)
	}

	word := d.token.Word
	if word == nil {
		return ip.basicRuntimeverr(format, args...)
	}

	traceLines := errors.Trace(word)

	lines := make([]string, 0, 8 + len(traceLines))
	lines = append(lines, "\n")
	lines = append(lines, fmt.Sprintf(format, args...))
	for _, line := range traceLines {
		lines = append(lines, fmt.Sprintf(" %v", line))
	}

	var sb strings.Builder
	for i, line := range lines {
		if i > 2 {
			sb.WriteString("\n")
		}

		sb.WriteString(line)
	}

	return &StepResult{
		status: ip.runtimev(format, args...),
		result: fmt.Errorf("%v\n", sb.String()),
	}
}

func typecheckv(format string, args...any) {
	if !IsVerboseTypeCheck() || len(format) <= 0 {
		return
	}

	colorYellow := "\033[33m"
	colorReset := "\033[0m"

	fmt.Printf("%sTypeCheck:%s ", colorYellow, colorReset)
	fmt.Printf(format, args...)

	if format[len(format)-1] != '\n' {
		fmt.Println()
	}
}
