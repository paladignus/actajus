package valueobject

import (
	"regexp"
	"testing"
)

type File string

func (f File) Value() string {
	return string(f)
}

func (f File) IsValid() bool {
	re := regexp.MustCompile(`^[\w,\s-]+\.(png|jpg|gif)$`)
	return re.MatchString(f.Value())
}

func TestFile(t *testing.T) {
	invalidFiles := []File{
		"document.pdf",     // Extensão inválida
		"image.jpeg",       // Extensão inválida
		"archive.zip",      // Extensão inválida
		"noextensionfile",  // Sem extensão
		"wrong.extension.", // Extensão inválida
		".hiddenfile",      // Sem nome antes do ponto
	}
	t.Run("should return false if file is invalid", func(t *testing.T) {
		for _, file := range invalidFiles {
			if file.IsValid() {
				t.Errorf("Expected File to be invalid, but got valid")
			}
		}
	})
}
