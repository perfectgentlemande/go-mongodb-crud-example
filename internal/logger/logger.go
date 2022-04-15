package logger

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrorField = "error"

type loggerCtxKey struct{}

func DefaultLogger() *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), os.Stdout, zap.DebugLevel))
}
func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, log)
}
func GetLogger(ctx context.Context) *zap.Logger {
	le, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger)
	if !ok {
		le = DefaultLogger()
	}
	return le
}
func NewLoggingMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextLog := log.With(
				zapcore.Field{
					Key:    "method",
					Type:   zapcore.StringType,
					String: r.Method,
				},
				zapcore.Field{
					Key:    "path",
					Type:   zapcore.StringType,
					String: r.URL.Path,
				})

			handler.ServeHTTP(w, r.WithContext(WithLogger(r.Context(), nextLog)))
		})
	}
}
