package tools

import "fmt"

// The browser build has no disk or network. The error gives the script the
// tool's normal failure result, and the interpreter warns about it on stderr.
func unavailable(tool string) error {
	return fmt.Errorf("%s is not available in the playground", tool)
}
