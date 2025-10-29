package otel

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/jvllmr/frans/internal/config"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otelsdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

func buildHTTPLogExporter(cfg config.OtelExclusive) (*otlploghttp.Exporter, error) {
	endpoint, err := url.Parse(cfg.HTTP.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.HTTP.Endpoint, err)
	}

	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(endpoint.Host),
		otlploghttp.WithHeaders(cfg.HTTP.Headers),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}

	if endpoint.Scheme == "http" {
		opts = append(opts, otlploghttp.WithInsecure())
	}

	if endpoint.Path != "" {
		opts = append(opts, otlploghttp.WithURLPath(endpoint.Path))
	}

	if cfg.HTTP.TLS.IsSet() {
		tlsConfig, err := cfg.HTTP.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlploghttp.WithTLSClientConfig(tlsConfig))
	}

	return otlploghttp.New(context.Background(), opts...)
}

func buildGRPCLogExporter(cfg config.OtelExclusive) (*otlploggrpc.Exporter, error) {
	host, port, err := net.SplitHostPort(cfg.GRPC.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.GRPC.Endpoint, err)
	}

	opts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(fmt.Sprintf("%s:%s", host, port)),
		otlploggrpc.WithHeaders(cfg.GRPC.Headers),
		otlploggrpc.WithCompressor(gzip.Name),
	}

	if cfg.GRPC.Insecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}

	if cfg.GRPC.TLS.IsSet() {
		tlsConfig, err := cfg.GRPC.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlploggrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	}

	return otlploggrpc.New(context.Background(), opts...)
}

func NewLoggerProvider(
	ctx context.Context,
	cfg config.OtelExclusive,
) (*otelsdk.LoggerProvider, error) {
	var (
		err      error
		exporter otelsdk.Exporter
	)
	if cfg.GRPC.Endpoint != "" {
		exporter, err = buildGRPCLogExporter(cfg)
	} else {
		exporter, err = buildHTTPLogExporter(cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("setting up exporter: %w", err)
	}

	var resAttrs []attribute.KeyValue
	for k, v := range cfg.ResourceAttributes {
		resAttrs = append(resAttrs, attribute.String(k, v))
	}

	res, err := resource.New(ctx,
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),

		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(config.FransVersion),
		),
		resource.WithAttributes(resAttrs...),

		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("building resource: %w", err)
	}

	bp := otelsdk.NewBatchProcessor(exporter)
	loggerProvider := otelsdk.NewLoggerProvider(
		otelsdk.WithResource(res),
		otelsdk.WithProcessor(bp),
	)

	return loggerProvider, nil
}
