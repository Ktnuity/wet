package tools

import (
	"fmt"
	"strings"
)

func ToolToStringInt(value int) string {
	return fmt.Sprintf("%d", value)
}

func ToolToStringString(value string) string {
	return value
}

func ToolToStringPath(path string) string {
	if strings.HasPrefix(path, ":") || strings.HasPrefix(path, "/") {
		return path[1:]
	} else if strings.HasPrefix(path, "./") {
		return path[2:]
	} else {
		return ""
	}
}
