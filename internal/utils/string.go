package utils

import "strings"

func LastSplit(s string, sep string) string {
	idx := strings.LastIndex(s, sep)
	if idx == -1 {
		return s
	}
	return s[idx+len(sep):]
}
