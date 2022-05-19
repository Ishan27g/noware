package noop

import (
	"context"
	"testing"

	"github.com/Ishan27g/noware/pkg/actions"
	"github.com/stretchr/testify/assert"
)

var serverEvent = actions.Event{
	Name: "server",
	Meta: nil,
}
var clientEvent = actions.Event{
	Name: "client",
	Meta: nil,
}
var noopCtx = func(ctx context.Context) context.Context {
	return context.WithValue(ctx, noopKey, true)
}

func TestNewCtxWithNoop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "no context", args: struct {
			ctx context.Context
		}{ctx: nil}, want: context.WithValue(context.Background(), noopKey, true)},
		{name: "some context", args: struct {
			ctx context.Context
		}{ctx: context.WithValue(context.WithValue(context.Background(), "one", "1"), "two", "2")},
			want: noopCtx(context.WithValue(context.WithValue(context.Background(), "one", "1"), "two", "2"))},
		{name: "noop context", args: struct {
			ctx context.Context
		}{ctx: noopCtx(context.Background())},
			want: noopCtx(context.Background())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewCtxWithNoop(tt.args.ctx), "NewCtxWithNoop(%v, %v)", tt.args.ctx)
		})
	}
}
