package web

import (
	"context"
	"time"
)

type key int8

const ctxKey key = 0

type Context struct {
	Now time.Time // Now is normalised to UTC
}

func setContext(ctx context.Context) context.Context {
	now := time.Now().UTC()
	ctx = context.WithValue(ctx, ctxKey, Context{
		Now: now,
	})
	return ctx
}

func GetContext(ctx context.Context) Context {
	v, ok := ctx.Value(ctxKey).(Context)
	if !ok {
		panic("invalid context value")
	}

	return v
}

// Now returns a time.Time that has been normalised to UTC. This value should be
// used for any and all time related activity to keep the time across a single
// request consistent.
func Now(ctx context.Context) time.Time {
	v, ok := ctx.Value(ctxKey).(Context)
	if !ok {
		panic("invalid context value")
	}

	return v.Now
}
