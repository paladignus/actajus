package module_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoCrossModuleLayerImports(t *testing.T) {
	t.Helper()

	var files []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatal("module directory not found")
		}
		t.Fatalf("walk module files: %v", err)
	}

	var violations []string
	for _, file := range files {
		srcModule, ok := moduleNameFromFile(file)
		if !ok {
			continue
		}

		fset := token.NewFileSet()
		parsed, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse imports for %s: %v", file, err)
		}

		for _, imp := range parsed.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			dstModule, layer, ok := moduleLayerFromImport(importPath)
			if !ok || dstModule == srcModule {
				continue
			}
			violations = append(violations, file+" imports "+importPath+" ("+layer+")")
		}
	}

	if len(violations) > 0 {
		t.Fatalf("cross-module imports into application/presentation/infrastructure are forbidden:\n%s", strings.Join(violations, "\n"))
	}
}

func moduleNameFromFile(path string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := 0; i+2 < len(parts); i++ {
		if parts[i] == "internal" && parts[i+1] == "module" {
			return parts[i+2], true
		}
	}
	return "", false
}

func moduleLayerFromImport(path string) (module string, layer string, ok bool) {
	parts := strings.Split(path, "/")
	for i := 0; i+4 < len(parts); i++ {
		if parts[i] == "internal" && parts[i+1] == "module" {
			layer = parts[i+3]
			if layer == "application" || layer == "presentation" || layer == "infrastructure" {
				return parts[i+2], layer, true
			}
			return "", "", false
		}
	}
	return "", "", false
}
