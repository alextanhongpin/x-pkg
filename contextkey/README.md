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

	return MustBe[T](val)
}

var UserContext Key[*User] = "user"

func main() {
	ctx := context.Background()
	ctx = UserContext.WithValue(ctx, &User{Name: "john"})
	res, ok := UserContext.Value(ctx)
	// res := UserContext.MustValue(ctx)
	fmt.Println(res, ok)

	fmt.Println(MustBe[string](10))
}

type User struct {
	Name string
}

func As[T any](unk any) (t T, ok bool) {
	t, ok = unk.(T)

	return
}

func MustBe[T any](unk any) T {
	t, ok := As[T](unk)
	if !ok {
		panic(fmt.Errorf("cast error: expected %T, got %T", t, unk))
	}

	return t
}
```
