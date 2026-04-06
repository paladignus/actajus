// Package validation
package validation

import "strings"

func parseRuleRaw(rule string) (name, raw string) {
	parts := strings.SplitN(rule, "=", 2)
	name = parts[0]
	if len(parts) == 2 {
		raw = parts[1]
	}
	return
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseRequiredIf(raw string) (field, expected string) {
	parts := strings.SplitN(raw, ",", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func parseJSONFieldName(tag, fallback string) string {
	if tag == "" || tag == "-" {
		return strings.ToLower(fallback)
	}
	name, _, _ := strings.Cut(tag, ",")
	name = strings.TrimSpace(name)
	if name == "" {
		return strings.ToLower(fallback)
	}
	return name
}
