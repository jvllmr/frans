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
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

func buildHTTPMetricExporter(cfg config.Otel) (*otlpmetrichttp.Exporter, error) {
	endpoint, err := url.Parse(cfg.HTTP.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.HTTP.Endpoint, err)
	}

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint.Host),
		otlpmetrichttp.WithHeaders(cfg.HTTP.Headers),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	}

	if endpoint.Scheme == "http" {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if endpoint.Path != "" {
		opts = append(opts, otlpmetrichttp.WithURLPath(endpoint.Path))
	}

	if cfg.HTTP.TLS.IsSet() {
		tlsConfig, err := cfg.HTTP.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlpmetrichttp.WithTLSClientConfig(tlsConfig))
	}

	return otlpmetrichttp.New(context.Background(), opts...)
}

func buildGRPCMetricExporter(cfg config.Otel) (*otlpmetricgrpc.Exporter, error) {
	host, port, err := net.SplitHostPort(cfg.GRPC.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid collector endpoint %q: %w", cfg.GRPC.Endpoint, err)
	}

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(fmt.Sprintf("%s:%s", host, port)),
		otlpmetricgrpc.WithHeaders(cfg.GRPC.Headers),
		otlpmetricgrpc.WithCompressor(gzip.Name),
	}

	if cfg.GRPC.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	if cfg.GRPC.TLS.IsSet() {
		tlsConfig, err := cfg.GRPC.TLS.CreateTLSConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating TLS client config: %w", err)
		}

		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)))
	}

	return otlpmetricgrpc.New(context.Background(), opts...)
}

func NewMeterProvider(
	ctx context.Context,
	cfg config.Otel,
) (func(), error) {
	var (
		err      error
		exporter otelsdk.Exporter
	)
	if cfg.GRPC.Endpoint != "" {
		exporter, err = buildGRPCMetricExporter(cfg)
	} else {
		exporter, err = buildHTTPMetricExporter(cfg)
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

	pr := otelsdk.NewPeriodicReader(
		exporter,
		otelsdk.WithTimeout(time.Second*time.Duration(cfg.PushInterval/2)),
		otelsdk.WithInterval(time.Second*time.Duration(cfg.PushInterval)),
	)
	mp := otelsdk.NewMeterProvider(
		otelsdk.WithResource(res),
		otelsdk.WithReader(pr))
	otel.SetMeterProvider(mp)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("Shutting down meter provider")
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down meter provider: %v", err)
		}
	}, nil

}
