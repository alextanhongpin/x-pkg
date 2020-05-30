package gojwt_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"errors"

	"github.com/alextanhongpin/pkg/gojwt"
	"github.com/dgrijalva/jwt-go"
)

func Example() {
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

func TestSignAndVerify(t *testing.T) {
	var (
		audience     = "your company"
		issuer       = "your service"
		secret       = "secret"
		subject      = "user 1"
		expiresAfter = 10 * time.Second
	)
	signer := gojwt.New(
		gojwt.Option{
			Secret:       []byte(secret),
			ExpiresAfter: expiresAfter,
			DefaultClaims: &gojwt.Claims{
				StandardClaims: jwt.StandardClaims{
					Audience: audience,
					Issuer:   issuer,
				},
			},
		},
	)
	token, err := signer.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = subject
		return nil
	})
	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}
	claims, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("verify token failed: %v", err)
	}

	if subject != claims.StandardClaims.Subject {
		t.Fatalf("expected %s, got %s", subject, claims.StandardClaims.Subject)
	}
}

func TestDifferentAudience(t *testing.T) {
	signer1 := gojwt.New(
		gojwt.Option{
			Secret:       []byte("secret"),
			ExpiresAfter: 10 * time.Second,
			DefaultClaims: &gojwt.Claims{
				StandardClaims: jwt.StandardClaims{
					Audience: "audience 1",
					Issuer:   "issuer",
				},
			},
		},
	)

	signer2 := gojwt.New(
		gojwt.Option{
			Secret:       []byte("secret"),
			ExpiresAfter: 10 * time.Second,
			DefaultClaims: &gojwt.Claims{
				StandardClaims: jwt.StandardClaims{
					Audience: "audience 2",
					Issuer:   "issuer",
				},
			},
			Validator: func(c *gojwt.Claims) error {
				if !c.StandardClaims.VerifyAudience("audience 2", true) {
					return fmt.Errorf("audience is invalid: %s", c.StandardClaims.Audience)
				}
				return nil
			},
		},
	)
	token, err := signer1.Sign(func(c *gojwt.Claims) error {
		c.StandardClaims.Subject = "user 1"
		return nil
	})
	if err != nil {
		t.Fatalf("signing failed: %v", err)
	}
	_, err = signer2.Verify(token)
	if err.Error() != "audience is invalid: audience 1" {
		t.Fatalf("expected %q, got %q", "audience is invalid: audience 1", err.Error())
	}
}
