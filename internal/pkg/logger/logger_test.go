package logger_test

import (
	"context"
	"testing"

	"github.com/hhk7734/gin-test/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithFields(t *testing.T) {
	cases := []struct {
		in   []zap.Field
		want []zap.Field
	}{
		{
			in: []zap.Field{
				zap.String("a", "a"),
				zap.Bool("b", true),
			},
			want: []zap.Field{
				zap.String("a", "a"),
				zap.Bool("b", true),
			},
		},
		{
			in: []zap.Field{
				zap.String("a", "a"),
				zap.String("a", "aa"),
			},
			want: []zap.Field{
				zap.String("a", "aa"),
			},
		},
		{
			in: []zap.Field{
				zap.String("a", "a"),
				zap.Reflect("a", zap.String("a", "a")),
			},
			want: []zap.Field{
				zap.Reflect("a", zap.String("a", "a")),
			},
		},
	}

	for _, c := range cases {
		ctx := context.Background()
		ctx = logger.WithFields(ctx, c.in...)
		assert.ElementsMatch(t, c.want, logger.Fields(ctx))
	}

	for _, c := range cases {
		ctx := context.Background()
		for _, f := range c.in {
			ctx = logger.WithFields(ctx, f)
		}
		assert.ElementsMatch(t, c.want, logger.Fields(ctx))
	}
}
