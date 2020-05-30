# gojwt

Usage:

```go
package main

import (
	"errors"
	"log"
	"time"

	"github.com/alextanhongpin/pkg/gojwt"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	var (
		audience     = "your company"
		issuer       = "your service"
		secret       = "secret"
		expiresAfter = 10 * time.Second
	)
	opt := gojwt.Option{
		Secret:       []byte(secret),
		ExpiresAfter: expiresAfter,
		DefaultClaims: &gojwt.Claims{
			StandardClaims: jwt.StandardClaims{
				Audience: audience,
				Issuer:   issuer,
			},
		},
		Validator: func(c *gojwt.Claims) error {
			if c.Issuer != issuer || c.Audience != audience {
				return errors.New("invalid token")
			}
			return nil
		},
	}
	signer := gojwt.New(opt)
	token, err := signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = "user id"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(token)
	claims, err := signer.Verify(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(claims)
}
```
