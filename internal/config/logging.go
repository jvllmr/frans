package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type LogConfig struct {
	LogJSON bool `mapstructure:"json"`
}

func setLogConfigDefaults(viper *viper.Viper) {
	viper.SetDefault("log.json", false)
}

type LogExclusive struct {
	LogConfig `mapstructure:"log"`
}

func NewLogConfig() (LogExclusive, error) {
	var cfg LogExclusive
	logConf := viper.New()
	setLogConfigDefaults(logConf)
	setConfigSearchStrategy(logConf)
	if err := logConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	if err := logConf.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
