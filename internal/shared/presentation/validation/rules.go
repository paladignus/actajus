// Package validation
package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func (v *Validator) applyRules(field string, value any, tag string) {
	if _, exists := v.errors[field]; exists {
		return
	}
	for _, rule := range strings.Split(tag, "|") {
		if _, exists := v.errors[field]; exists {
			return
		}
		name, param := parseRule(rule)
		switch name {
		case "required":
			if isZero(value) {
				v.addError(field, fmt.Sprintf(messages[v.lang]["required"], field))
			}
		case "min":
			validateMin(v, field, value, param)
		case "len":
			validateLen(v, field, value, param)
		case "numeric":
			validateNumeric(v, field, value)
		}
	}
}

func validateMin(v *Validator, field string, value any, min int) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		if len(rv.String()) < min {
			v.addError(field, fmt.Sprintf(messages[v.lang]["min"], field, min))
		}
	case reflect.Int:
		if rv.Int() < int64(min) {
			v.addError(field, fmt.Sprintf(messages[v.lang]["min"], field, min))
		}
	}
}

func validateLen(v *Validator, field string, value any, length int) {
	s := fmt.Sprintf("%v", value)
	if len(s) != length {
		v.addError(field, fmt.Sprintf(messages[v.lang]["len"], field, length))
	}
}

func validateNumeric(v *Validator, field string, value any) {
	s := fmt.Sprintf("%v", value)
	clean := strings.NewReplacer("-", "", ".", "", "/", "").Replace(s)
	matched, _ := regexp.MatchString(`^\d+$`, clean)
	if !matched {
		v.addError(field, fmt.Sprintf(messages[v.lang]["numeric"], field))
	}
}
