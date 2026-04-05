package orchestrator

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func NewTraceContext(ctx context.Context) context.Context{
	newCtx := context.WithValue(ctx, "traceID", uuid.New().String())
	return newCtx
}

func GetLogger(ctx context.Context, source string) *slog.Logger{
	traceID := ctx.Value("traceID").(string)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(
		[]slog.Attr{slog.String("traceID", traceID), slog.String("source", source)},
	),
	)
	return logger
}