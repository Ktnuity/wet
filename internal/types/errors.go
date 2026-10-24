package types

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	Err error
}

type GenericError struct {
	Err error
}

func (e *RuntimeError) Error() string {
	colorRed := "\033[31m"
	colorReset := "\033[0m"

	var sb strings.Builder
	sb.WriteString(colorRed)
	sb.WriteString("Runtime Error: ")
	sb.WriteString(colorReset)
	sb.WriteString("%v")
	return fmt.Sprintf(sb.String(), e.Err)
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("%v\n", e.Err)
}

func MakeError(err error) error {
	return &GenericError{ Err: err }
}
