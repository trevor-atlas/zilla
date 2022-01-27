package util

import "os"

// Exists returns whether the given file or directory exists or not
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
