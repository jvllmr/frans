package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/otel"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
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

func SetupLogging() error {
	logCfg, err := config.NewLogConfig()
	if err != nil {
		return fmt.Errorf("get log config: %w", err)
	}

	var stdoutHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if logCfg.LogJSON {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	otelCfg, err := config.NewOtelConfig()

	if err != nil {
		return fmt.Errorf("get otel config: %w", err)
	}

	loggingProvider, err := otel.NewLoggerProvider(context.Background(), otelCfg)

	if err != nil {
		return fmt.Errorf("setup otel logging: %v", err)
	}

	logger := slog.New(
		slogmulti.Fanout(
			otelslog.NewHandler("frans", otelslog.WithLoggerProvider(loggingProvider)),
			stdoutHandler,
		),
	)

	slog.SetDefault(logger)
	return nil
}
