package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DBType     string `mapstructure:"db_type"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     uint16 `mapstructure:"db_port"`
	DBName     string `mapstructure:"db_name"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
}

func setDBConfigDefaults(viper *viper.Viper) {
	viper.SetDefault("db_type", "postgres")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 0)
	viper.SetDefault("db_name", "frans")
	viper.SetDefault("db_user", "frans")
	viper.SetDefault("db_password", "")
}

func NewDBConfig() (DBConfig, error) {
	var config DBConfig
	dbConf := viper.New()
	setDBConfigDefaults(dbConf)
	setConfigSearchStrategy(dbConf)
	if err := dbConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	if err := dbConf.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
