package contextkey_test

import (
	"context"
	"testing"

	"github.com/alextanhongpin/pkg/contextkey"
)

func TestContextKey(t *testing.T) {
	var (
		userIDKey = contextkey.Key("user_id")
		id        = "1234"
	)
	ctx := userIDKey.WithValue(context.Background(), id)
	got, ok := userIDKey.Value(ctx).(string)
	if true != ok {
		t.Fatalf("expected %t, got %t", true, ok)
	}
	if id != got {
		t.Fatalf("expected %s, got %s", id, got)
	}
}
