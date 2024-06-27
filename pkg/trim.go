package pkg

import (
	"strings"
)

func TrimData(b []byte) string {
	d := strings.ReplaceAll(string(b), "\n", "")
	d = strings.ReplaceAll(d, "\t", "")
	d = strings.ReplaceAll(d, "\\", "")
	d = strings.ReplaceAll(d, "\"", "")
	return d
}
