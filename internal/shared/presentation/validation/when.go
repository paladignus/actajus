// Package validation
package validation

func When(cond bool, fn func()) {
	if cond {
		fn()
	}
}
