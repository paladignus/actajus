// Package validation
package validation

import (
	"reflect"
	"time"
)

type Validator struct {
	shortCircuitGlobal bool
}

type Option func(*Validator)

func New(opts ...Option) *Validator {
	v := &Validator{}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *Validator) ValidateStruct(input any) []Violation {
	var out []Violation
	v.validateStruct(input, "", &out)
	return out
}

func (v *Validator) validateStruct(input any, prefix string, out *[]Violation) {
	if v.shortCircuitGlobal && len(*out) > 0 {
		return
	}
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	if !val.IsValid() {
		return
	}
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}
		val = val.Elem()
		typ = typ.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < val.NumField(); i++ {
		if v.shortCircuitGlobal && len(*out) > 0 {
			return
		}
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		if fieldType.PkgPath != "" {
			continue
		}
		field := parseJSONFieldName(fieldType.Tag.Get("json"), fieldType.Name)
		path := field
		if prefix != "" {
			path = prefix + "." + path
		}
		if tag := fieldType.Tag.Get("validate"); tag != "" {
			applyRules(path, input, fieldVal.Interface(), tag, out)
		}
		if fieldVal.Kind() == reflect.Struct && fieldVal.Type() != reflect.TypeOf(time.Time{}) {
			v.validateStruct(fieldVal.Interface(), path, out)
		}
		if fieldVal.Kind() == reflect.Slice {
			for j := 0; j < fieldVal.Len(); j++ {
				item := fieldVal.Index(j)
				if item.Kind() == reflect.Struct {
					v.validateStruct(item.Interface(), path+"["+itoa(j)+"]", out)
				}
			}
		}
	}
}
