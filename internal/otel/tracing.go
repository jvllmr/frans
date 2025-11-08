package otel

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/url"
	"time"

	"github.com/jvllmr/frans/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const TracingService = "frans"

func buildHTTPTracingExporter(cfg config.Otel) (*otlptrace.Exporter, error) {
	endpoint, err := url.Parse(cfg.HTTP.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.HTTP.Endpoint, err)
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint.Host),
		otlptracehttp.WithHeaders(cfg.HTTP.Headers),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	}

	if endpoint.Scheme == "http" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if endpoint.Path != "" {
		opts = append(opts, otlptracehttp.WithURLPath(endpoint.Path))
	}

	if cfg.HTTP.TLS.IsSet() {
		tlsConfig, err := cfg.HTTP.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlptracehttp.WithTLSClientConfig(tlsConfig))
	}

	return otlptrace.New(context.Background(), otlptracehttp.NewClient(opts...))
}

func buildGRPCTracingExporter(cfg config.Otel) (*otlptrace.Exporter, error) {
	host, port, err := net.SplitHostPort(cfg.GRPC.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.GRPC.Endpoint, err)
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", host, port)),
		otlptracegrpc.WithHeaders(cfg.GRPC.Headers),
		otlptracegrpc.WithCompressor(gzip.Name),
	}

	if cfg.GRPC.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	if cfg.GRPC.TLS.IsSet() {
		tlsConfig, err := cfg.GRPC.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	}

	return otlptrace.New(context.Background(), otlptracegrpc.NewClient(opts...))
}

func NewTracerProvider(ctx context.Context, cfg config.Otel) (func(), error) {
	slog.Info("Setup open telemetry tracing provider...")
	var (
		err      error
		exporter *otlptrace.Exporter
	)
	if cfg.GRPC.Endpoint != "" {
		exporter, err = buildGRPCTracingExporter(cfg)
	} else {
		exporter, err = buildHTTPTracingExporter(cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("setting up tracing exporter: %w", err)
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

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("Shutting down tracing provider")
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}, err
}

func GetFransTracer() trace.Tracer {
	return otel.Tracer(TracingService)
}

func NewSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tr := GetFransTracer()
	return tr.Start(ctx, spanName)
}
