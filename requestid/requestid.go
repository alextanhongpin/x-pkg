package requestid

import (
	"context"
	"net/http"
)

// XRequestID represents the header constant. NOTE: the Id is not uppercase.
const XRequestID = "X-Request-Id"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const requestID = contextKey(XRequestID)

// Provider represents the interface for the request id provider.
type Provider func(http.ResponseWriter, *http.Request) (string, error)

// Factory represents the function to generate a request id.
type Factory func() (string, error)

// RequestID reads the request id from the current request header, and creates
// a new one if it does not exist.
func RequestID(factory Factory) Provider {
	return func(w http.ResponseWriter, r *http.Request) (string, error) {
		id := r.Header.Get(XRequestID)
		if id == "" {
			var err error
			id, err = factory()
			if err != nil {
				return "", err
			}
			r.Header.Set(XRequestID, id)
		}
		w.Header().Set(XRequestID, id)
		*r = *r.WithContext(WithValue(r.Context(), id))
		return id, nil
	}
}

// WithValue populates the context with the given request id.
func WithValue(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestID, reqID)
}

// Value extracts the request id from the provided context.
func Value(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestID).(string)
	return id, ok
}
