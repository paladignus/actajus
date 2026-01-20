package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RepositorySpy struct {
	user IUser
	log  Logger
}

func (r *RepositorySpy) User() IUser {
	return r.user
}

func (r *RepositorySpy) Logger() Logger {
	return r.log
}

func TestRepositoryInterface(t *testing.T) {
	sut := &RepositorySpy{
		user: &UserSpy{},
		log:  &LoggerSpy{},
	}
	user := sut.User()
	assert.NotNil(t, user)
	log := sut.Logger()
	assert.NotNil(t, log)
}
