package config

import (
	"fmt"
	"net/http"
	"strings"

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

type Config struct {
	DBConfig `mapstructure:",squash"`
	DevMode  bool `mapstructure:"dev_mode"`

	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	RootPath string `mapstructure:"root_path"`

	OidcIssuer       string `mapstructure:"oidc_issuer"`
	OidcClientID     string `mapstructure:"oidc_client_id"`
	OidcClientSecret string `mapstructure:"oidc_client_secret"`
	OidcAdminGroup   string `mapstructure:"oidc_admin_group"`

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

func setConfigSearchStrategy(viper *viper.Viper) {
	viper.SetConfigName("frans")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("frans")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func setDBConfigDefaults(viper *viper.Viper) {
	viper.SetDefault("db_type", "postgres")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 0)
	viper.SetDefault("db_name", "frans")
	viper.SetDefault("db_user", "frans")
	viper.SetDefault("db_password", "")
}

func GetDBConfig() (DBConfig, error) {
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

	setDBConfigDefaults(fransConf)

	fransConf.SetDefault("oidc_admin_group", "admin")

	fransConf.SetDefault("log_json", false)

	fransConf.SetDefault("smtp_port", 25)
	fransConf.SetDefault("smtp_username", nil)
	fransConf.SetDefault("smtp_password", nil)

	setConfigSearchStrategy(fransConf)

	if err := fransConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

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
