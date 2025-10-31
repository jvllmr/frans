package logging

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// adapted from https://github.com/FabienMht/ginslog
func RecoveryLogger(logger *slog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {

				attributes := []slog.Attr{}

				attributes = append(attributes, slog.Any("error", r))

				attributes = append(attributes, slog.String("stack", string(debug.Stack())))

				// Log the panic
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					"Panic recovered",
					attributes...)

				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}
				span := trace.SpanFromContext(c.Request.Context())
				span.RecordError(err, trace.WithStackTrace(true))
				span.SetStatus(codes.Error, "internal error")
			}
		}()
		c.Next()
	}
}
