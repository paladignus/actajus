package valueobject

import (
	"testing"

	valueobject "github.com/paladignus/actajus/internal/domain/value_object"
)

func TestFile(t *testing.T) {
	invalidFiles := []valueobject.File{
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
	validFiles := []valueobject.File{
		"image.png",
		"document.pdf",
		"photo.jpg",
		"graphic.gif",
		"picture.jpeg",
	}
	t.Run("should return true if file is valid", func(t *testing.T) {
		for _, file := range validFiles {
			if !file.IsValid() {
				t.Errorf("Expected File to be valid, but got invalid")
			}
		}
	})
}
