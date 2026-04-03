package query

import (
	"strconv"
	"strings"
)

func OptionalInt16(value string) int16 {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 16)
	if err != nil || parsed <= 0 {
		return 0
	}
	return int16(parsed)
}
