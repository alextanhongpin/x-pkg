package gojwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/pkg/errors"
)

// package main
//
// import (
//         "errors"
//         "log"
//         "time"
//
//         "github.com/alextanhongpin/pkg/gojwt"
//         "github.com/dgrijalva/jwt-go"
// )
//
// func main() {
//         var (
//                 audience     = "your company"
//                 issuer       = "your service"
//                 semver       = "0.0.1"
//                 scope        = "guest"
//                 role         = "user"
//                 secret       = "secret"
//                 expiresAfter = 10 * time.Second
//         )
//         opt := gojwt.Option{
//                 Secret:       []byte(secret),
//                 ExpiresAfter: expiresAfter,
//                 DefaultClaims: &gojwt.Claims{
//                         Semver: semver,
//                         Scope:  scope,
//                         Role:   role,
//                         StandardClaims: jwt.StandardClaims{
//                                 Audience: audience,
//                                 Issuer:   issuer,
//                         },
//                 },
//                 Validator: func(c *gojwt.Claims) error {
//                         if c.Semver != semver ||
//                                 c.Issuer != issuer ||
//                                 c.Audience != audience {
//                                 return errors.New("invalid token")
//                         }
//                         return nil
//                 },
//         }
//         signer := gojwt.New(opt)
//         token, err := signer.Sign(func(c *gojwt.Claims) error {
//                 c.StandardClaims.Subject = "user id"
//                 return nil
//         })
//         if err != nil {
//                 log.Fatal(err)
//         }
//         log.Println(token)
//         claims, err := signer.Verify(token)
//         if err != nil {
//                 log.Fatal(err)
//         }
//         log.Println(claims)
// }

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
