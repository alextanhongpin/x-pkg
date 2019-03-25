package gojwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/pkg/errors"
)

type (
	Claims struct {
		Semver string
		Scope  string
		Role   string
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
		return "", errors.Wrap(err, "sign token failed")
	}
	// Set the expires at and issued at time.
	claims.ExpiresAt = now.Add(expiresAfter).Unix()
	claims.IssuedAt = now.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	return ss, errors.Wrap(err, "sign token failed")
}

// Verify checks if the given token string is valid, and returns the claims or
// error.
func (j *JwtSigner) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.opt.Secret, nil
		},
	)
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
