package logging

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// adapted from https://github.com/FabienMht/ginslog
func RecoveryLogger(logger *slog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				attributes := []slog.Attr{}

				attributes = append(attributes, slog.Any("error", err))

				attributes = append(attributes, slog.String("stack", string(debug.Stack())))

				// Log the panic
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					"Panic recovered",
					attributes...)

			}
		}()
		c.Next()
	}
}
