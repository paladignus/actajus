// Package validation
package validation

func (v *Validator) When(cond bool, fn func()) {
	if cond {
		fn()
	}
}
