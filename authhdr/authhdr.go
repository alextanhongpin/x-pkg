// Package authhdr reduces the boilerplate code required to extract authentication
// token from the Authorization header.
package authhdr

import (
	"errors"
	"net/http"
	"strings"
)

const (
	Bearer = "bearer"
	Basic  = "basic"
)

// AuthHeader represents the entity for the authorization header.
type AuthHeader struct {
	bearer string
	token  string
}

// New returns a new AuthHeader
func New() *AuthHeader {
	return new(AuthHeader)
}

// Extract attempts to obtain the authorization bearer and token from the `Authorization` header.
func (a *AuthHeader) Extract(r *http.Request) error {
	auth := r.Header.Get("Authorization")
	paths := strings.Fields(auth)
	if len(paths) != 2 {
		return errors.New("invalid authorization header")
	}
	a.bearer, a.token = paths[0], paths[1]
	return nil
}

// BearerIs checks if the given bearer is equal to the compared bearer.
func (a *AuthHeader) BearerIs(s string) bool {
	return strings.EqualFold(a.bearer, s)
}

// Bearer returns the authorization bearer.
func (a *AuthHeader) Bearer() string {
	return a.bearer
}

// Token returns the authorization token.
func (a *AuthHeader) Token() string {
	return a.token
}
