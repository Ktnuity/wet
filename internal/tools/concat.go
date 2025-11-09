package tools

import (
	"fmt"
	"os"
)

func ToolConcatPath(res, str string) string {
	return fmt.Sprintf("%s%c%s", res, os.PathSeparator, str)
}

func ToolConcatString(a, b string) string {
	return fmt.Sprintf("%s%s", a, b)
}
