// Package validation
package validation

import (
	"reflect"
	"strings"
)

func isZero(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Int, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + (n % 10))
		n /= 10
	}
	return string(b[i:])
}

func hasPath(vs []Violation, path string) bool {
	for i := range vs {
		if vs[i].Path == path {
			return true
		}
	}
	return false
}
