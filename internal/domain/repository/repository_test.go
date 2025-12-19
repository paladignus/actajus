package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RepositorySpy struct {
	user  IUser
	log   Logger
	cache Cache
}

func (r *RepositorySpy) User() IUser {
	return r.user
}

func (r *RepositorySpy) Logger() Logger {
	return r.log
}

func (r *RepositorySpy) Cache() Cache {
	return r.cache
}

func TestRepositoryInterface(t *testing.T) {
	sut := &RepositorySpy{
		// user:  &AuthenticationSpy{},
		log:   &LoggerSpy{},
		cache: &CacheSpy{},
	}
	user := sut.User()
	assert.NotNil(t, user)
	log := sut.Logger()
	assert.NotNil(t, log)
	cache := sut.Cache()
	assert.NotNil(t, cache)
}
