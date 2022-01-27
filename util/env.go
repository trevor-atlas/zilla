package util

import (
	"os"
	"strings"
)

func getenv() map[string]string {
	m := make(map[string]string)

	for _, e := range os.Environ() {
		if strings.Contains(e, "ZILLA") {
			pair := strings.SplitN(e, "=", 2)
			m[pair[0]] = pair[1]
		}
	}
	return m
}
