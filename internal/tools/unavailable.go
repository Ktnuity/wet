//go:build !js

package tools

// Every tool is available outside the browser build.
func unavailable(tool string) error {
	return nil
}
