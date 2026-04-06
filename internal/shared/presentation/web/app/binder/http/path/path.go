package pathbinder

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func BindInt64Path(r *http.Request, key string) (int64, error) {
	return parseInt64Value(r.PathValue(key), "invalid path id")
}

func BindInt16Path(r *http.Request, key string) (int16, error) {
	return parseInt16Value(r.PathValue(key), "invalid path id")
}

func ParseInt64Value(raw string, message string) (int64, error) {
	return parseInt64Value(raw, message)
}

func ParseInt16Value(raw string, message string) (int16, error) {
	return parseInt16Value(raw, message)
}

func parseInt64Value(raw string, message string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value <= 0 {
		return 0, errors.New(message)
	}
	return value, nil
}

func parseInt16Value(raw string, message string) (int16, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 16)
	if err != nil || value <= 0 {
		return 0, errors.New(message)
	}
	return int16(value), nil
}
