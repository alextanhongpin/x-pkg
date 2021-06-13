package gojwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	UserInfo struct {
		Email      string                 `json:"email,omitempty"`
		Verified   bool                   `json:"verified,omitempty"`
		FamilyName string                 `json:"family_name,omitempty"`
		GivenName  string                 `json:"given_name,omitempty"`
		Locale     string                 `json:"locale,omitempty"`
		Name       string                 `json:"name,omitempty"`
		Picture    string                 `json:"picture,omitempty"`
		Extra      map[string]interface{} `json:"extra,omitempty"`
	}

	Claims struct {
		UserInfo
		jwt.StandardClaims
	}

	Validator func(*Claims) error

	NowFunc func() time.Time

	// Option can be configured based on your needs.
	Option struct {
		Secret        []byte
		ExpiresAfter  time.Duration
		DefaultClaims *Claims
		Validator     Validator
		NowFunc       NowFunc
	}

	// Signer represents the JwtSigner operations.
	Signer interface {
		Sign(func(*Claims) error) (string, error)
		Verify(token string) (*Claims, error)
	}

	// JwtSigner ...
	JwtSigner struct {
		opt Option
	}
)

var ErrEmpty = errors.New("Authorization header is not provided")

// DefaultNowFunc returns the time of signing the token. Provides an interface
// for mocking the time.
func DefaultNowFunc() time.Time {
	return time.Now()
}

// DefaultValidator return a Nop validator.
func DefaultValidator(*Claims) error {
	return nil
}

// New returns a new jwt signer.
func New(opt Option) *JwtSigner {
	if opt.DefaultClaims == nil {
		opt.DefaultClaims = &Claims{}
	}
	if opt.Validator == nil {
		opt.Validator = DefaultValidator
	}
	if opt.NowFunc == nil {
		opt.NowFunc = DefaultNowFunc
	}
	return &JwtSigner{opt}
}

// Sign takes a function that modifies the claims and return a signed token
// string.
func (j *JwtSigner) Sign(fn func(c *Claims) error) (string, error) {
	var (
		claims       = *j.opt.DefaultClaims
		expiresAfter = j.opt.ExpiresAfter
		now          = j.opt.NowFunc()
		secret       = j.opt.Secret
	)
	err := fn(&claims)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}
	// Set the expires at and issued at time.
	claims.ExpiresAt = now.Add(expiresAfter).Unix()
	claims.IssuedAt = now.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}
	return ss, nil
}

// Verify checks if the given token string is valid, and returns the claims or
// error.
func (j *JwtSigner) Verify(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrEmpty
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.opt.Secret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid authorization header: %w", err)
	}

	// Apparently this is possible by sending Authorization: Bearer
	// undefined.
	if token == nil {
		return nil, errors.New("invalid authorization header")
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, j.opt.Validator(claims)
	}
	return nil, err
}
