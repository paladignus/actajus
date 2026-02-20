// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type AccessClaims struct {
	idUser    vo.ID
	idSession vo.ID
	roles     []string
	issuer    string
	audience  string
	expiresAt time.Time
}

type AccessClaimsBuilder struct {
	claims *AccessClaims
}

func NewAccessClaimsBuilder() *AccessClaimsBuilder {
	return &AccessClaimsBuilder{
		claims: &AccessClaims{},
	}
}

func (a *AccessClaimsBuilder) WithIDUser(id int64) *AccessClaims {
	a.claims.idUser = vo.ID(id)
	return a.claims
}

func (a *AccessClaimsBuilder) WithIDSession(id int64) *AccessClaims {
	a.claims.idSession = vo.ID(id)
	return a.claims
}

func (a *AccessClaimsBuilder) WithRoles(roles []string) *AccessClaims {
	a.claims.roles = roles
	return a.claims
}

func (a *AccessClaimsBuilder) WithIssuer(issuer string) *AccessClaims {
	a.claims.issuer = issuer
	return a.claims
}

func (a *AccessClaimsBuilder) WithAudience(audience string) *AccessClaims {
	a.claims.audience = audience
	return a.claims
}

func (a *AccessClaimsBuilder) WithExpiresAt(expiresAt time.Time) *AccessClaims {
	a.claims.expiresAt = expiresAt
	return a.claims
}

func (a *AccessClaimsBuilder) Build() (*AccessClaims, error) {
	return a.claims, nil
}

func (a *AccessClaims) IDUser() vo.ID        { return a.idUser }
func (a *AccessClaims) IDSession() vo.ID     { return a.idSession }
func (a *AccessClaims) Roles() []string      { return a.roles }
func (a *AccessClaims) Issuer() string       { return a.issuer }
func (a *AccessClaims) Audience() string     { return a.audience }
func (a *AccessClaims) ExpiresAt() time.Time { return a.expiresAt }
