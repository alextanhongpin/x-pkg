package main

// https://console.bluemix.net/docs/services/cloud-object-storage/hmac/hmac-signature.html#constructing-an-hmac-signature
// https://dev.to/mathijspim/hmac-authentication-better-protection-for-your-api-4e0
// https://docs.aws.amazon.com/AWSECommerceService/latest/DG/HMACSignatures.html

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	Header struct {
		name, value string
	}

	Option struct {
		Bearer       string
		ExpiresAfter time.Duration
	}

	Repository interface {
		LookupSecretKey(accessKeyID string) (string, error)
	}

	Signer interface {
		ConvertMapToHeaders(fields map[string]interface{}) http.Header
		SignHeaders(secretKey string, header http.Header) string
		ValidateHeaderDate(header http.Header) error
		ValidateHeader(header http.Header) error
		NewAuthorizationHeader(accessKeyID, signature string) string
		SelectHeadersWithPrefix(header http.Header) http.Header
	}

	SignerImpl struct {
		opt          Option
		headerPrefix string
		headerDate   string
		repo         Repository
	}
)

func NewSigner(opt Option, repo Repository) *SignerImpl {
	opt.Bearer = strings.Title(opt.Bearer)
	return &SignerImpl{
		opt:          opt,
		headerPrefix: newHeaderPrefix(opt.Bearer),
		headerDate:   newHeaderDate(opt.Bearer),
		repo:         repo,
	}
}

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

func (s *SignerImpl) ValidateHeaderDate(header http.Header) error {
	var date time.Time
	dateStr := header.Get(s.headerDate)
	dateInt, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		return err
	}
	date = time.Unix(dateInt, 0)
	if time.Since(date) > s.opt.ExpiresAfter {
		return errors.New("token expired")
	}
	return nil
}

func (s *SignerImpl) NewAuthorizationHeader(accessKeyID, signature string) string {
	return fmt.Sprintf("%s %s:%s", s.opt.Bearer, accessKeyID, signature)
}

func (s *SignerImpl) ValidateHeader(header http.Header) error {
	authorization := header.Get("Authorization")
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 {
		// Error message should be conveying a message - x is invalid, x is required
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
	if err := s.ValidateHeaderDate(header); err != nil {
		return err
	}
	headersWithPrefix := s.SelectHeadersWithPrefix(header)
	encodedSignature := s.SignHeaders(secretKey, headersWithPrefix)
	if subtle.ConstantTimeCompare([]byte(encodedSignature), []byte(signature)) != 1 {
		return errors.New("invalid signature")
	}
	return nil
}

func (s *SignerImpl) SelectHeadersWithPrefix(header http.Header) http.Header {
	result := make(http.Header)
	l := len(s.headerPrefix)
	for key, values := range header {
		if strings.EqualFold(key[:min(l, len(key))], s.headerPrefix) {
			result.Set(key, values[0])
		}
	}
	return result
}

type repository struct {
}

func (r *repository) LookupSecretKey(accessKeyID string) (string, error) {
	if accessKeyID == "xyz" {
		return "secret", nil
	}
	return "", errors.New("not found")
}

func main() {
	var (
		secretKey    = "secret"
		accessKeyID  = "xyz"
		bearer       = "custom"
		expiresAfter = 5 * time.Second
	)

	repo := &repository{}
	opt := Option{
		Bearer:       bearer,
		ExpiresAfter: expiresAfter,
	}
	signer := NewSigner(opt, repo)

	fields := map[string]interface{}{
		"key9": "Value9",
		"keyA": "Value2",
		"keyZ": "Value3",
		"key1": "Value1",
		// "date": time.Now().Add(-10 * time.Second).Unix(),
		// "X-John-Date": time.Now().Add(-10 * time.Second).Unix(),
	}

	header := signer.ConvertMapToHeaders(fields)
	fmt.Println("header", header)
	signature := signer.SignHeaders(secretKey, header)

	// TODO: Get query from url querystring too.
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header = header
	req.Header.Add("Authorization", signer.NewAuthorizationHeader(accessKeyID, signature))
	req.Header.Add("X-Request-ID", "xyz")
	fmt.Println(req.Header)

	// http.Header and url.Values can be conversible.
	cvt := url.Values(req.Header)
	fmt.Println(http.Header(cvt))

	err := signer.ValidateHeader(req.Header)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signature is valid")
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
