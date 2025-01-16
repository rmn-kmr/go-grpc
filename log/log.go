package log

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/metadata"
	"os"
)

const CHAR_LIMIT = 100
const REQUEST_ID = "x-request-id"

func Info(ctx context.Context, msg string, fields ...any) {
	msg, fields = buildLogFields(ctx, msg, fields)
	slog.Info(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...any) {
	msg, fields = buildLogFields(ctx, msg, fields)
	slog.Debug(msg, fields...)
}

func Error(ctx context.Context, err error, msg string, fields ...any) {

	msg, fields = buildLogFields(ctx, msg, fields)
	// add error to fields
	fields = append(fields, "error", err)
	slog.Error(msg, fields...)
}

func buildLogFields(ctx context.Context, msg string, fields []any) (string, []any) {
	requestID := middleware.GetReqID(ctx)
	if requestID == "" {
		if values := metadata.ValueFromIncomingContext(ctx, REQUEST_ID); len(values) > 0 {
			requestID = values[0]
		}
	}
	fields = append(fields, "request_id", requestID)

	// limit message to CHAR_LIMIT
	if len(msg) > CHAR_LIMIT {
		msg = msg[:CHAR_LIMIT] + "..."
	}
	return msg, fields
}

func InitLogger(logLevel string) {
	// set log level for slog
	var level slog.Level
	// parse level from string
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
