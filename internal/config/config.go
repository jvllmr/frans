package config

import (
	"fmt"
	"net/http"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	DevMode bool `mapstructure:"dev_mode"`

	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	RootPath string `mapstructure:"root_path"`

	OidcIssuer       string `mapstructure:"oidc_issuer"`
	OidcClientID     string `mapstructure:"oidc_client_id"`
	OidcClientSecret string `mapstructure:"oidc_client_secret"`
	OidcAdminGroup   string `mapstructure:"oidc_admin_group"`

	DBType     string `mapstructure:"db_type"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     uint16 `mapstructure:"db_port"`
	DBName     string `mapstructure:"db_name"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`

	FilesDir string `mapstructure:"files_dir"`
	MaxSizes int64  `mapstructure:"max_sizes"`
	MaxFiles uint8  `mapstructure:"max_files"`

	DefaultExpiryDaysSinceLastDownload uint8 `mapstructure:"expiry_days_since"`
	DefaultExpiryTotalDownloads        uint8 `mapstructure:"expiry_total_dl"`
	DefaultExpiryTotalDays             uint8 `mapstructure:"expiry_total_days"`

	GrantDefaultExpiryDaysSinceLastUpload uint8 `mapstructure:"grant_expiry_days_since"`
	GrantDefaultExpiryTotalUploads        uint8 `mapstructure:"grant_expiry_total_up"`
	GrantDefaultExpiryTotalDays           uint8 `mapstructure:"grant_expiry_total_days"`

	SMTPServer   string  `mapstructure:"smtp_server"`
	SMTPPort     int     `mapstructure:"smtp_port"`
	SMTPFrom     string  `mapstructure:"smtp_from"`
	SMTPUsername *string `mapstructure:"smtp_username"`
	SMTPPassword *string `mapstructure:"smtp_password"`

	LogJSON bool `mapstructure:"log_json"`
}

func GetConfig() (Config, error) {
	var config Config
	fransConf := viper.New()

	fransConf.SetDefault("dev_mode", false)
	fransConf.SetDefault("root_path", "")

	fransConf.SetDefault("host", "127.0.0.1")
	fransConf.SetDefault("port", 8080)

	fransConf.SetDefault("files_dir", "files")
	fransConf.SetDefault("max_files", 20)
	fransConf.SetDefault("max_sizes", 2_000_000_000) // 2GB

	fransConf.SetDefault("expiry_days_since", 7)
	fransConf.SetDefault("expiry_total_dl", 10)
	fransConf.SetDefault("expiry_total_days", 30)

	fransConf.SetDefault("grant_expiry_days_since", 7)
	fransConf.SetDefault("grant_expiry_total_up", 10)
	fransConf.SetDefault("grant_expiry_total_days", 30)

	fransConf.SetDefault("db_type", "postgres")
	fransConf.SetDefault("db_host", "localhost")
	fransConf.SetDefault("db_port", 0)
	fransConf.SetDefault("db_name", "frans")
	fransConf.SetDefault("db_user", "frans")
	fransConf.SetDefault("db_password", "")

	fransConf.SetDefault("oidc_admin_group", "admin")

	fransConf.SetDefault("log_json", false)

	fransConf.SetDefault("smtp_port", 25)
	fransConf.SetDefault("smtp_username", nil)
	fransConf.SetDefault("smtp_password", nil)

	fransConf.SetConfigName("frans")
	fransConf.SetConfigType("yaml")
	fransConf.AddConfigPath(".")
	fransConf.SetEnvPrefix("frans")
	fransConf.AutomaticEnv()
	fransConf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := fransConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	// Unmarshal into the struct
	if err := fransConf.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}

func GetSafeConfig() Config {
	configValue, err := GetConfig()
	if err != nil {
		panic(err)
	}
	return configValue
}

func GetBaseURL(configValue Config, request *http.Request) string {
	proto := "http"
	if request.TLS != nil {
		proto = "https"
	}
	host := request.Host
	patchedRootPath := configValue.RootPath
	if len(patchedRootPath) == 0 {
		patchedRootPath = "/"
	}
	return fmt.Sprintf("%s://%s%s", proto, host, patchedRootPath)
}
