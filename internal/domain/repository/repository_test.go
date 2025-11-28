package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RepositorySpy struct {
	auth  Authentication
	log   Logger
	cache Cache
}

func (r *RepositorySpy) Authentication() Authentication {
	return r.auth
}

func (r *RepositorySpy) Logger() Logger {
	return r.log
}

func (r *RepositorySpy) Cache() Cache {
	return r.cache
}

func TestRepositoryInterface(t *testing.T) {
	sut := &RepositorySpy{
		auth:  &AuthenticationSpy{},
		log:   &LoggerSpy{},
		cache: &CacheSpy{},
	}
	auth := sut.Authentication()
	assert.NotNil(t, auth)
	log := sut.Logger()
	assert.NotNil(t, log)
	cache := sut.Cache()
	assert.NotNil(t, cache)
}

