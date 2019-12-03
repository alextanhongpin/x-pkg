package contextkey_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alextanhongpin/pkg/contextkey"
)

func Example() {
	var (
		userIDKey = contextkey.Key("user_id")
		id        = "1234"
	)
	ctx := userIDKey.WithValue(context.Background(), id)
	got, ok := userIDKey.Value(ctx).(string)
	fmt.Println(got, ok)
	// Output: 1234, true
}

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
