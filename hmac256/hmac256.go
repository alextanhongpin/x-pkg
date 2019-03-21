package main

// https://console.bluemix.net/docs/services/cloud-object-storage/hmac/hmac-signature.html#constructing-an-hmac-signature
// https://dev.to/mathijspim/hmac-authentication-better-protection-for-your-api-4e0
// https://docs.aws.amazon.com/AWSECommerceService/latest/DG/HMACSignatures.html

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const BEARER = "John"
const HEADER_PREFIX = "X-John-"
const HEADER_DATE = "X-John-Date"
const EXPIRES_AFTER = 5 * time.Second

func main() {
	var secretKey = "secret"
	fields := map[string]interface{}{
		"key9": "Value9",
		"keyA": "Value2",
		"keyZ": "Value3",
		"key1": "Value1",
	}
	headerHeader := headersFromMap(fields)
	// Generate signature
	headerSignature := signHeaders(secretKey, headerHeader)

	// TODO: Get query from url querystring too.
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header = headerHeader
	req.Header.Add("Authorization", newAuthorizationHeader("xyz", headerSignature))
	req.Header.Add("X-Request-ID", "xyz")
	fmt.Println(req.Header)
	header := req.Header

	authHeader := req.Header.Get("Authorization")
	paths := strings.Split(authHeader, " ")
	if len(paths) != 2 {
		log.Fatal("authorization header is invalid")
	}
	bearer, token := paths[0], paths[1]
	if bearer != BEARER {
		log.Fatalf(`bearer "%s" is invalid`, bearer)
	}
	credentials := strings.Split(token, ":")
	if len(credentials) != 2 {
		log.Fatal("authorization token is invalid")
	}
	accessKeyID, signature := credentials[0], credentials[1]
	fmt.Println("accessKeyID=", accessKeyID, "signature=", signature)

	var ts time.Time
	{
		s := req.Header.Get(HEADER_DATE)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ts = time.Unix(i, 0)
		fmt.Println("date", s, i, ts)
	}
	fmt.Println("has expired?", time.Since(ts) > EXPIRES_AFTER, time.Since(ts))

	headerWithPrefix := selectHeaderPrefix(header)
	// Lookup secret.
	var secret = ""
	if accessKeyID == "xyz" {
		secret = "secret"
	}
	encodedSignature := signHeaders(secret, headerWithPrefix)
	fmt.Println(encodedSignature == signature)

}

func newAuthorizationHeader(accessKeyID, signature string) string {
	return fmt.Sprintf("%s %s:%s", BEARER, accessKeyID, signature)
}
func selectHeaderPrefix(header http.Header) http.Header {
	newHeader := make(http.Header)
	for key, values := range header {
		if yes := strings.HasPrefix(key, HEADER_PREFIX); yes && len(values) > 0 {
			newHeader.Set(key, values[0])
		}
	}
	return newHeader
}
func headersFromMap(m map[string]interface{}) http.Header {
	header := make(http.Header, 0)
	for k, v := range m {
		if !strings.HasPrefix(strings.ToLower(k), strings.ToLower(HEADER_PREFIX)) {
			k = HEADER_PREFIX + k
		}
		header.Set(k, fmt.Sprint(v))
	}
	return header
}

func signHeaders(secretKey string, header http.Header) string {
	// Set if not exist - mandatory field.
	if v := header.Get(HEADER_DATE); v == "" {
		header.Set(HEADER_DATE, strconv.FormatInt(time.Now().Unix(), 10))
	}
	fmt.Println("generated header", header)
	var headers []struct {
		name, value string
	}
	for key, values := range header {
		headers = append(headers, struct{ name, value string }{key, values[0]})
	}
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].name < headers[j].name
	})
	result := make([]string, len(headers))
	for i, h := range headers {
		result[i] = fmt.Sprintf("%s:%s", h.name, h.value)
	}
	return createSignature(secretKey, strings.Join(result, " "))
}

func createSignature(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// UseCase: Verify Header
// Select authorization header
// Split into bearer and token
// Check if bearer is equals your custom bearer
// Split token into accessKeyID and signature
// Get headers
// Take headers with the prefix x-yourcustomheader-field and its value
// Check if the header x-yourcustomheader-date is present.
// Sort the fields in ascending order
// Concatenate the values and compute the signature with the accessKeyID
// Compare the computed signature and the signature in header

// UseCase: Generate Signature
// Populate map
// Add prefix for missing name
// Convert map to http header 
// Sign headers with secret key
// Set authorization header
// Fire call
