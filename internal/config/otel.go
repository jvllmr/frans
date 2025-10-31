package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type ClientTLS struct {
	CA                 string `mapstructure:"ca,omitempty"`
	Cert               string `mapstructure:"cert,omitempty"`
	Key                string `mapstructure:"key,omitempty"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify,omitempty"`
}

func (c *ClientTLS) IsSet() bool {
	return c != nil && (c.CA != "" || c.Cert != "" || c.Key != "" || c.InsecureSkipVerify != false)
}

func (c *ClientTLS) CreateTLSConfig(ctx context.Context) (*tls.Config, error) {
	if c == nil {
		return nil, nil
	}

	var caPool *x509.CertPool

	if c.CA != "" {
		var ca []byte
		if _, errCA := os.Stat(c.CA); errCA == nil {
			var err error
			ca, err = os.ReadFile(c.CA)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA. %w", err)
			}
		} else {
			ca = []byte(c.CA)
		}

		caPool = x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(ca) {
			return nil, errors.New("failed to parse CA")
		}
	}

	hasCert := len(c.Cert) > 0
	hasKey := len(c.Key) > 0

	if hasCert != hasKey {
		return nil, errors.New("both TLS cert and key must be defined")
	}

	if !hasCert || !hasKey {
		return &tls.Config{
			RootCAs:            caPool,
			InsecureSkipVerify: c.InsecureSkipVerify,
		}, nil
	}

	cert, err := loadKeyPair(c.Cert, c.Key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: c.InsecureSkipVerify,
	}, nil
}

func loadKeyPair(cert, key string) (tls.Certificate, error) {
	keyPair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err == nil {
		return keyPair, nil
	}

	_, err = os.Stat(cert)
	if err != nil {
		return tls.Certificate{}, errors.New("cert file does not exist")
	}

	_, err = os.Stat(key)
	if err != nil {
		return tls.Certificate{}, errors.New("key file does not exist")
	}

	keyPair, err = tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return keyPair, nil
}

type OtelGRPC struct {
	Endpoint string            `mapstructure:"endpoint,omitempty"`
	Insecure bool              `mapstructure:"insecure,omitempty"`
	TLS      *ClientTLS        `mapstructure:"tls,omitempty"`
	Headers  map[string]string `mapstructure:"headers,omitempty"`
}

type OtelHTTP struct {
	Endpoint string            `mapstructure:"endpoint,omitempty"`
	TLS      *ClientTLS        `mapstructure:"tls,omitempty"`
	Headers  map[string]string `mapstructure:"headers,omitempty"`
}

type Otel struct {
	GRPC               *OtelGRPC         `mapstructure:"grpc,omitempty"`
	HTTP               *OtelHTTP         `mapstructure:"http,omitempty"`
	PushInterval       int               `mapstructure:"push_interval,omitempty"`
	ResourceAttributes map[string]string `mapstructure:"resource_attributes,omitempty"`
	ServiceName        string            `mapstructure:"service_name,omitempty"`
	SampleRate         float64           `mapstructure:"sample_rate,omitempty"`
}

type OtelExclusive struct {
	Otel `mapstructure:"otel"`
}

func setOtelConfigDefaults(viper *viper.Viper) {
	// general settings
	viper.SetDefault("otel.service_name", "frans")
	viper.SetDefault("otel.resource_attributes", map[string]string{})
	viper.SetDefault("otel.push_interval", 10)
	viper.SetDefault("otel.sample_rate", 1.0)

	// grpc
	viper.SetDefault("otel.grpc.endpoint", "")
	viper.SetDefault("otel.grpc.insecure", false)
	viper.SetDefault("otel.grpc.headers", map[string]string{})
	viper.SetDefault("otel.grpc.tls.ca", "")
	viper.SetDefault("otel.grpc.tls.cert", "")
	viper.SetDefault("otel.grpc.tls.key", "")
	viper.SetDefault("otel.grpc.tls.insecure_skip_verify", "")

	// http
	viper.SetDefault("otel.http.endpoint", "")
	viper.SetDefault("otel.http.headers", map[string]string{})
	viper.SetDefault("otel.http.tls.ca", "")
	viper.SetDefault("otel.http.tls.cert", "")
	viper.SetDefault("otel.http.tls.key", "")
	viper.SetDefault("otel.http.tls.insecure_skip_verify", "")

}

func NewOtelConfig() (OtelExclusive, error) {
	var cfg OtelExclusive
	otelConf := viper.New()
	setOtelConfigDefaults(otelConf)
	setConfigSearchStrategy(otelConf)
	if err := otelConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	if err := otelConf.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode into struct: %w", err)
	}
	return cfg, nil
}
