package contextkey

import "context"

type Key string

func (k Key) Value(ctx context.Context) interface{} {
	return ctx.Value(k)
}

func (k Key) WithValue(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, k, value)
}
