package usecase

import (
	"testing"

	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/spy"
)

type SignIn struct {
	repository repository.AuthenticateRepository
}

func NewSignIn(repository repository.AuthenticateRepository) SignIn {
	return SignIn{repository}
}

func (s SignIn) Execute(username, password string) (string, error) {
	return s.repository.Authenticate(username, password)
}

func TestSignIn(t *testing.T) {
	username := NewRandomString(10)
	password := NewRandomString(10)
	repository := spy.AuthenticateSpyPersistence{}
	sut := NewSignIn(&repository)
	sut.Execute(username, password)

	t.Run("should corrects properties from repository", func(t *testing.T) {
		if repository.Username != username || repository.Password != password {
			t.Error("Expected username and password to match")
		}
	})

	t.Run("should call repository only once", func(t *testing.T) {
		if repository.CallsCount != 1 {
			t.Error("Expected callsCount to be 1")
		}
	})

	t.Run("should return user id", func(t *testing.T) {
		userID, _ := sut.Execute(username, password)
		if userID == "" {
			t.Error("Expected userID to be not empty")
		}
	})

	t.Run("should return an error if username id empty", func(t *testing.T) {
		_, err := sut.Execute("", password)
		if err == nil {
			t.Error("Expected error to be nil")
		}
	})

	t.Run("should return an error if password id empty", func(t *testing.T) {
		_, err := sut.Execute(username, "")
		if err == nil {
			t.Error("Expected error to be nil")
		}
	})
}

func NewRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
