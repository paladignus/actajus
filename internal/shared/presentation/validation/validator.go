// Package validation
package validation

import (
	"reflect"
	"strings"

	"github.com/paladignus/actajus/internal/shared/domain"
)

type Validator struct {
	lang   Lang
	errors map[string]*domain.FieldError
}

func New(lang Lang) *Validator {
	return &Validator{
		lang:   lang,
		errors: make(map[string]*domain.FieldError),
	}
}

func (v *Validator) addError(field, message string) {
	if _, exists := v.errors[field]; exists {
		return
	}
	v.errors[field] = domain.NewFieldError(field, message)
}

func (v *Validator) ValidateStruct(input any) error {
	return v.validateStruct(input)
}

func (v *Validator) validateStruct(input any) error {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		field := fieldType.Tag.Get("json")
		if field == "" || field == "-" {
			field = strings.ToLower(fieldType.Name)
		}
		if tag := fieldType.Tag.Get("validate"); tag != "" {
			v.applyRules(field, fieldVal.Interface(), tag)
		}
		if fieldVal.Kind() == reflect.Struct && !isTime(fieldVal.Type()) {
			v.validateStruct(fieldVal.Interface())
		}
		if fieldVal.Kind() == reflect.Slice {
			for i := 0; i < fieldVal.Len(); i++ {
				item := fieldVal.Index(i)
				if item.Kind() == reflect.Struct {
					v.validateStruct(
						item.Interface(),
					)
				}
			}
		}
	}
	if len(v.errors) > 0 {
		return domain.NewValidationErrors(v.toSlice())
	}
	return nil
}

func isTime(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

//	func (v *Validator) ValidateStruct(input any) error {
//		val := reflect.ValueOf(input)
//		typ := reflect.TypeOf(input)
//		if val.Kind() == reflect.Pointer {
//			val = val.Elem()
//			typ = typ.Elem()
//		}
//		for i := 0; i < val.NumField(); i++ {
//			fieldVal := val.Field(i)
//			fieldType := typ.Field(i)
//			tag := fieldType.Tag.Get("validate")
//			if tag == "" {
//				continue
//			}
//			field := fieldType.Tag.Get("json")
//			if field == "" {
//				field = fieldType.Name
//			}
//			v.applyRules(field, fieldVal.Interface(), tag)
//		}
//
//		if len(v.errors) > 0 {
//			return sharedDomain.NewValidationErrors(v.toSlice())
//		}
//		return nil
//	}
func (v *Validator) toSlice() []*domain.FieldError {
	out := make([]*domain.FieldError, 0, len(v.errors))
	for _, e := range v.errors {
		out = append(out, e)
	}
	return out
}
