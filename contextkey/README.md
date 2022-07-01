# Contextkey

Easily setup context.

## go 1.18

With generics, the solution will be:
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
)

type Key[T any] string

func (c Key[T]) WithValue(ctx context.Context, t T) context.Context {
	return context.WithValue(ctx, c, t)
}

func (c Key[T]) Value(ctx context.Context) (t T, ok bool) {
	val := ctx.Value(c)
	if val == nil {
		return
	}

	t, ok = val.(T)

	return
}

func (c Key[T]) MustValue(ctx context.Context) T {
	val := ctx.Value(c)
	if val == nil {
		panic(fmt.Errorf("contextkey: not found: %s", c))
	}

	t, ok := val.(T)
	if !ok {
		panic(fmt.Errorf("contextkey: not found: %s", c))
	}

	return t
}

var UserContext Key[*User] = "user"

func main() {
	ctx := context.Background()
	ctx = UserContext.WithValue(ctx, &User{Name: "john"})
	res, ok := UserContext.Value(ctx)
	// res := UserContext.MustValue(ctx)
	fmt.Println(res, ok)
}

type User struct {
	Name string
}
```
