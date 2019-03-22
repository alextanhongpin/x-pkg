package hmac256

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	// Header wraps the http.Header name and value.
	Header struct {
		name, value string
	}

	// Option represents the configurable signer option.
	Option struct {
		Bearer       string
		ExpiresAfter time.Duration
	}

	// Repository represents a lookup interface for the secret key.
	Repository interface {
		LookupSecretKey(accessKeyID string) (string, error)
	}

	// Signer creates header and validates them.
	Signer interface {
		ConvertMapToHeaders(fields map[string]interface{}) http.Header
		SignHeaders(secretKey string, header http.Header) string
		ValidateHeader(header http.Header) error
		NewAuthorizationHeader(accessKeyID, signature string) string
	}

	// SignerImpl implements the Signer interface.
	SignerImpl struct {
		opt          Option
		repo         Repository
		headerDate   string
		headerPrefix string
	}
)

// NewSigner returns a new hmac 256 signer.
func NewSigner(opt Option, repo Repository) *SignerImpl {
	opt.Bearer = strings.Title(opt.Bearer)
	return &SignerImpl{
		opt:          opt,
		repo:         repo,
		headerDate:   newHeaderDate(opt.Bearer),
		headerPrefix: newHeaderPrefix(opt.Bearer),
	}
}

// ConvertMapToHeaders takes in a map, adds the additional prefix if not
// present, and returns a new http.Header.
func (s *SignerImpl) ConvertMapToHeaders(fields map[string]interface{}) http.Header {
	header := make(http.Header)
	h := len(s.headerPrefix)
	for k, v := range fields {
		if !strings.EqualFold(k[:min(h, len(k))], s.headerPrefix) {
			k = s.headerPrefix + k
		}
		header.Set(k, fmt.Sprint(v))
	}
	return header
}

// SignHeaders takes a http.Header, orders them alphabetically, and concatenate
// the header name and values before signing them with a secret key.
func (s *SignerImpl) SignHeaders(secretKey string, header http.Header) string {
	if v := header.Get(s.headerDate); v == "" {
		header.Set(s.headerDate, strconv.FormatInt(time.Now().Unix(), 10))
	}
	var headers []Header
	for key, values := range header {
		headers = append(headers, Header{key, values[0]})
	}
	return createSignature(secretKey, concatHeaders(headers...))
}

// NewAuthorizationHeader takes in the accessKeyID and signature and returns a
// new Authorization header.
func (s *SignerImpl) NewAuthorizationHeader(accessKeyID, signature string) string {
	return fmt.Sprintf("%s %s:%s", s.opt.Bearer, accessKeyID, signature)
}

// ValidateHeader takes a http.Header, checks if the Authorization header is
// valid and attempts to reconstruct the signature to check the validity of the
// payload.
func (s *SignerImpl) ValidateHeader(header http.Header) error {
	authorization := header.Get("Authorization")
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 {
		// Error message should be conveying a message - x is invalid,
		// x is required
		return errors.New("authorization header is invalid")
	}
	bearer, token := parts[0], parts[1]
	if bearer != s.opt.Bearer {
		return fmt.Errorf(`bearer "%s" is invalid`, bearer)
	}
	tokenParts := strings.Split(token, ":")
	if len(tokenParts) != 2 {
		return errors.New("token is invalid")
	}
	accessKey, signature := tokenParts[0], tokenParts[1]
	secretKey, err := s.repo.LookupSecretKey(accessKey)
	if err != nil {
		return err
	}
	if err := validateHeaderDate(header, s.headerDate, s.opt.ExpiresAfter); err != nil {
		return err
	}
	headersWithPrefix := selectHeadersWithPrefix(s.headerPrefix, header)
	encodedSignature := s.SignHeaders(secretKey, headersWithPrefix)
	if subtle.ConstantTimeCompare([]byte(encodedSignature), []byte(signature)) != 1 {
		return errors.New("invalid signature")
	}
	return nil
}

func validateHeaderDate(header http.Header, headerDate string, expiresAfter time.Duration) error {
	var date time.Time
	dateStr := header.Get(headerDate)
	dateInt, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		return err
	}
	date = time.Unix(dateInt, 0)
	if time.Since(date) > expiresAfter {
		return errors.New("token expired")
	}
	return nil
}

func selectHeadersWithPrefix(prefix string, header http.Header) http.Header {
	result := make(http.Header)
	l := len(prefix)
	for key, values := range header {
		if strings.EqualFold(key[:min(l, len(key))], prefix) {
			result.Set(key, values[0])
		}
	}
	return result
}

func newHeaderPrefix(bearer string) string {
	return fmt.Sprintf("X-%s-", bearer)
}

func newHeaderDate(bearer string) string {
	return fmt.Sprintf("X-%s-Date", bearer)
}

func concatHeaders(headers ...Header) string {
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].name < headers[j].name
	})
	result := make([]string, len(headers))
	for i, h := range headers {
		result[i] = fmt.Sprintf("%s:%s", h.name, h.value)
	}
	return strings.Join(result, " ")
}

func createSignature(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func min(hd int, rest ...int) int {
	for _, n := range rest {
		if n < hd {
			hd = n
		}
	}
	return hd
}
