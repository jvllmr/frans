package logging

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
)

// adapted from https://github.com/FabienMht/ginslog
func GinLogger(logger *slog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		c.Header("X-Request-ID", requestID)

		// Process the request
		c.Next()

		attributes := []slog.Attr{}

		attributes = append(attributes, slog.String("ip", c.ClientIP()))

		attributes = append(attributes, slog.Int("status", c.Writer.Status()))

		attributes = append(attributes, slog.String("method", c.Request.Method))

		attributes = append(attributes, slog.String("path", c.Request.URL.Path))

		attributes = append(attributes, slog.String("user-agent", c.Request.UserAgent()))

		attributes = append(attributes, slog.Duration("latency", time.Since(start)))

		attributes = append(attributes, slog.String("request-id", requestID))

		if len(c.Errors) > 0 {
			attributes = append(attributes, slog.String("errors", c.Errors.String()))
		}

		logger.LogAttrs(context.Background(), slog.LevelInfo, "Incoming request", attributes...)

	}
}

func SetupLogging(logConfig config.LogConfig) {
	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if logConfig.LogJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	basicLogger := slog.New(logHandler)
	slog.SetDefault(basicLogger)
}
