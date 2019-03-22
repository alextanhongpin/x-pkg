# hmac256

Usage:

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alextanhongpin/pkg/hmac256"
)

// repository is a non-concurrent safe data structure that holds the
// accessKeyID and secretKey pair. For production, consider a concurrent-safe
// implementation of map or use a persistent datastore.
type repository struct {
	data map[string]string
}

func (r *repository) LookupSecretKey(accessKeyID string) (string, error) {
	secretKey, exist := r.data[accessKeyID]
	if !exist {
		return "", errors.New("secret not found")
	}
	return secretKey, nil
}

var repo *repository

func init() {
	repo = &repository{
		data: map[string]string{
			"abc": "xyz",
		},
	}
}

func main() {
	var (
		accessKeyID  = "abc"
		bearer       = "custom"
		expiresAfter = 5 * time.Second
		secretKey    = "xyz"
	)

	opt := hmac256.Option{
		Bearer:       bearer,
		ExpiresAfter: expiresAfter,
	}
	signer := hmac256.NewSigner(opt, repo)

	fields := map[string]interface{}{
		"key9": "Value9",
		"keyA": "Value2",
		"keyZ": "Value3",
		"key1": "Value1",
		// If not provided, the date header will be automatically
		// populated.
		// "date": time.Now().Add(-10 * time.Second).Unix(),
		// "X-John-Date": time.Now().Add(-10 * time.Second).Unix(),
	}
	// Convert the map to a header, and sign the header to get the
	// signature.
	header := signer.ConvertMapToHeaders(fields)
	signature := signer.SignHeaders(secretKey, header)

	// Attach the signature in the Authorization header together with the
	// headers fields.
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header = header
	req.Header.Add("Authorization", signer.NewAuthorizationHeader(accessKeyID, signature))
	req.Header.Add("X-Request-ID", "xyz")
	// Make the call to the server...
	fmt.Println(req.Header)

	// http.Header and url.Values are conversible - which means we can pass
	// the values through the querystring too (e.g. for GET operations)
	// cvt := url.Values(http.Header{})
	// hdr := http.Header(url.Values{})

	// This can be part of a middleware that protects the endpoint that
	// needs the header to be validated.
	err := signer.ValidateHeader(req.Header)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signature is valid")
}
```

## Resources

- https://console.bluemix.net/docs/services/cloud-object-storage/hmac/hmac-signature.html#constructing-an-hmac-signature
- https://dev.to/mathijspim/hmac-authentication-better-protection-for-your-api-4e0
- https://docs.aws.amazon.com/AWSECommerceService/latest/DG/HMACSignatures.html
