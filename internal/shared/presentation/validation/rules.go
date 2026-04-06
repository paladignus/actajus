// Package validation
package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

func applyRules(path string, root, value any, tag string, out *[]Violation) {
	if hasPath(*out, path) {
		return
	}
	for rule := range strings.SplitSeq(tag, "|") {
		if hasPath(*out, path) {
			return
		}
		name, raw := parseRuleRaw(rule)
		switch name {
		case "required":
			if isZero(value) {
				*out = append(*out, Violation{Path: path, Code: CodeRequired})
			}
		case "min":
			validateMin(path, value, atoi(raw), out)
		case "len":
			validateLen(path, value, atoi(raw), out)
		case "numeric":
			validateNumeric(path, value, out)
		case "oneof":
			validateOneOf(path, value, splitCSV(raw), out)
		case "email":
			validateEmail(path, value, out)
		case "password":
			validatePassword(path, value, out)
		case "required_if":
			refField, expected := parseRequiredIf(raw)
			if refField != "" && evalRequiredIf(root, refField, expected) {
				if isZero(value) {
					*out = append(*out, Violation{
						Path: path,
						Code: CodeRequiredIf,
						Meta: map[string]string{
							"field":  refField,
							"equals": expected,
						},
					})
				}
			}
		}
	}
}

func validateMin(path string, value any, min int, out *[]Violation) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		if len(strings.TrimSpace(rv.String())) < min {
			*out = append(*out, Violation{Path: path, Code: CodeMin, Meta: map[string]string{"min": itoa(min)}})
		}
	case reflect.Int, reflect.Int64:
		if rv.Int() < int64(min) {
			*out = append(*out, Violation{Path: path, Code: CodeMin, Meta: map[string]string{"min": itoa(min)}})
		}
	}
}

func validateLen(path string, value any, length int, out *[]Violation) {
	s := fmt.Sprintf("%v", value)
	if len(s) != length {
		*out = append(*out, Violation{Path: path, Code: CodeLen, Meta: map[string]string{"len": itoa(length)}})
	}
}

func validateNumeric(path string, value any, out *[]Violation) {
	s := fmt.Sprintf("%v", value)
	clean := strings.NewReplacer("-", "", ".", "", "/", "").Replace(s)
	matched, _ := regexp.MatchString(`^\d+$`, clean)
	if !matched {
		*out = append(*out, Violation{Path: path, Code: CodeNumeric})
	}
}

func validateOneOf(path string, value any, allowed []string, out *[]Violation) {
	val := fmt.Sprintf("%v", value)
	if slices.Contains(allowed, val) {
		return
	}
	// for _, a := range allowed {
	// 	if val == a {
	// 		return
	// 	}
	// }
	*out = append(*out, Violation{
		Path: path,
		Code: CodeOneOf,
		Meta: map[string]string{
			"allowed": strings.Join(allowed, ","),
		},
	})
}

func validateEmail(path string, value any, out *[]Violation) {
	s := fmt.Sprintf("%v", value)
	if strings.TrimSpace(s) == "" {
		return
	}
	if _, err := mail.ParseAddress(s); err != nil {
		*out = append(*out, Violation{Path: path, Code: CodeEmail})
	}
}

func validatePassword(path string, value any, out *[]Violation) {
	s := fmt.Sprintf("%v", value)
	if strings.TrimSpace(s) == "" {
		return
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(s)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(s)
	if hasLower && hasUpper && hasNumber && hasSpecial {
		return
	}
	*out = append(*out, Violation{
		Path: path,
		Code: CodePassword,
		Meta: map[string]string{
			"requires": "lower,upper,number,special",
		},
	})
}

func evalRequiredIf(root any, refField, expected string) bool {
	rv := reflect.ValueOf(root)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		jsonName := parseJSONFieldName(sf.Tag.Get("json"), sf.Name)
		if jsonName != refField && sf.Name != refField {
			continue
		}
		fv := rv.Field(i).Interface()
		return fmt.Sprintf("%v", fv) == expected
	}
	return false
}
