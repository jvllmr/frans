package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type LogConfig struct {
	LogJSON bool `mapstructure:"log_json"`
}

func setLogConfigDefaults(viper *viper.Viper) {
	viper.SetDefault("log_json", false)
}

func NewLogConfig() (LogConfig, error) {
	var config LogConfig
	logConf := viper.New()
	setDBConfigDefaults(logConf)
	setConfigSearchStrategy(logConf)
	if err := logConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	if err := logConf.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
