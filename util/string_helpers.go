package util

import "strings"

func IsString(v any) bool {
	_, ok := v.(string)
	return ok
}

func IsStringEmpty(v any) bool {
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}
