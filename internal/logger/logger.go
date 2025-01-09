package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

var setupOnce sync.Once

var logger *zap.Logger

func Get() *zap.Logger {
	setupOnce.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		level := zap.NewAtomicLevelAt(zap.InfoLevel)

		cfg := zap.NewProductionEncoderConfig()
		cfg.TimeKey = "timestamp"
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(cfg)

		core := zapcore.NewCore(consoleEncoder, stdout, level)
		logger = zap.New(core)
	})
	return logger
}

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}
	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}
	return context.WithValue(ctx, contextKey{}, l)
}
