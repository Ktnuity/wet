package util

func AsRef[T any](value T) *T {
	return &value
}
