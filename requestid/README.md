## RequestID Package

Usage:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/alextanhongpin/pkg/requestid"
	"github.com/rs/xid"
)

type middleware func(next http.HandlerFunc) http.HandlerFunc

func withRequestID(provider requestid.Provider, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := provider(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}
func main() {
	provider := requestid.RequestID(func() (string, error) {
		// Provide your own implementation to generate the request id.
		// return "xyz", nil
		return xid.New().String(), nil
	})
	http.Handle("/", withRequestID(provider, index))
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	reqid, exist := requestid.Value(r.Context())
	fmt.Fprintf(w, "got request id: %s %t", reqid, exist)
}
```
