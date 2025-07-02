package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DevMode          bool   `mapstructure:"dev_mode"`
	RootPath         string `mapstructure:"root_path"`
	Host             string `mapstructure:"host"`
	Port             uint16 `mapstructure:"port"`
	OidcIssuer       string `mapstructure:"oidc_issuer"`
	OidcClientID     string `mapstructure:"oidc_client_id"`
	OidcClientSecret string `mapstructure:"oidc_client_secret"`
	DBType           string `mapstructure:"db_type"`
	DBHost           string `mapstructure:"db_host"`
	DBPort           uint16 `mapstructure:"db_port"`
	DBName           string `mapstructure:"db_name"`
	DBUser           string `mapstructure:"db_user"`
	DBPassword       string `mapstructure:"db_password"`
	AdminGroup       string `mapstructure:"admin_group"`
	MaxFiles         string `mapstructure:"max_files"`
	MaxSizes         string `mapstructure:"max_sizes"`
	FilesDir         string `mapstructure:"files_dir"`

	DefaultExpiryDaysSinceLastDownload uint8 `mapstructure:"default_expiry_days_since_last_download"`
	DefaultExpiryTotalDownloads        uint8 `mapstructure:"default_expiry_total_downloads"`
	DefaultExpiryTotalDays             uint8 `mapstructure:"default_expiry_total_days"`
}

func GetConfig() (Config, error) {
	var config Config
	fransConf := viper.New()

	fransConf.SetDefault("dev_mode", false)
	fransConf.SetDefault("root_path", "")

	fransConf.SetDefault("host", "127.0.0.1")
	fransConf.SetDefault("port", 8080)

	fransConf.SetDefault("max_files", 20)
	fransConf.SetDefault("max_sizes", 2_000_000_000) // 2GB

	fransConf.SetDefault("default_expiry_days_since_last_download", 7)
	fransConf.SetDefault("default_expiry_total_downloads", 10)
	fransConf.SetDefault("default_expiry_total_days", 30)

	fransConf.SetDefault("db_type", "postgres")
	fransConf.SetDefault("db_host", "localhost")
	fransConf.SetDefault("db_port", 0)
	fransConf.SetDefault("db_name", "frans")
	fransConf.SetDefault("db_user", "frans")
	fransConf.SetDefault("db_password", "")

	fransConf.SetDefault("admin_group", "admin")

	fransConf.SetConfigName("frans")
	fransConf.SetConfigType("yaml")
	fransConf.AddConfigPath(".")
	fransConf.SetEnvPrefix("frans")
	fransConf.AutomaticEnv()
	fransConf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := fransConf.ReadInConfig(); err != nil {
		log.Println("Warning: No config file found, falling back to environment variables.")
	}

	// Unmarshal into the struct
	if err := fransConf.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
