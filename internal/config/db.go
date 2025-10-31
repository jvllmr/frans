package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DBType     string `mapstructure:"type"`
	DBHost     string `mapstructure:"host"`
	DBPort     uint16 `mapstructure:"port"`
	DBName     string `mapstructure:"name"`
	DBUser     string `mapstructure:"user"`
	DBPassword string `mapstructure:"password"`
}

type DBExclusive struct {
	DBConfig `mapstructure:"db"`
}

func setDBConfigDefaults(viper *viper.Viper) {
	viper.SetDefault("db.type", "postgres")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 0)
	viper.SetDefault("db.name", "frans")
	viper.SetDefault("db.user", "frans")
	viper.SetDefault("db.password", "")
}

func NewDBConfig() (DBExclusive, error) {
	var config DBExclusive
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
