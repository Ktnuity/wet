package tools

import (
	"fmt"
)

func ToolToToken(str string) string {
	return fmt.Sprintf(":%s", str)
}

func ToolToAbsolute(str string) string {
	return fmt.Sprintf("/%s", str)
}

func ToolToRelative(str string) string {
	return fmt.Sprintf("./%s", str)
}
