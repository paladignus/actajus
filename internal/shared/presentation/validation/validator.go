// Package validation
package validation

import (
	"reflect"

	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type Validator struct {
	lang   Lang
	errors map[string]*sharedDomain.FieldError
}

func New(lang Lang) *Validator {
	return &Validator{
		lang:   lang,
		errors: make(map[string]*sharedDomain.FieldError),
	}
}

func (v *Validator) addError(field, message string) {
	if _, exists := v.errors[field]; exists {
		return
	}
	v.errors[field] = sharedDomain.NewFieldError(field, message)
}

func (v *Validator) ValidateStruct(input any) error {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}
		field := fieldType.Tag.Get("json")
		if field == "" {
			field = fieldType.Name
		}
		v.applyRules(field, fieldVal.Interface(), tag)
	}

	if len(v.errors) > 0 {
		return sharedDomain.NewValidationErrors(v.toSlice())
	}
	return nil
}

func (v *Validator) toSlice() []*sharedDomain.FieldError {
	out := make([]*sharedDomain.FieldError, 0, len(v.errors))
	for _, e := range v.errors {
		out = append(out, e)
	}
	return out
}
