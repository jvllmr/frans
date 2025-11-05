package config

import (
	"fmt"
	"net/http"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
)

type OidcConfig struct {
	OidcIssuer     string `mapstructure:"issuer"`
	OidcClientID   string `mapstructure:"client_id"`
	OidcAdminGroup string `mapstructure:"admin_group"`
}

type SMTPConfig struct {
	SMTPServer   string  `mapstructure:"server"`
	SMTPPort     int     `mapstructure:"port"`
	SMTPFrom     string  `mapstructure:"from"`
	SMTPUsername *string `mapstructure:"username"`
	SMTPPassword *string `mapstructure:"password"`
}

type FilesConfig struct {
	FilesDir string `mapstructure:"dir"`
	MaxSizes int64  `mapstructure:"max_size"`
	MaxFiles uint8  `mapstructure:"max_per_upload"`
}

type ExpiryConfig struct {
	DefaultExpiryDaysSinceLastDownload uint8 `mapstructure:"days_since_last_download"`
	DefaultExpiryTotalDownloads        uint8 `mapstructure:"total_downloads"`
	DefaultExpiryTotalDays             uint8 `mapstructure:"total_days"`
}

type GrantExpiryConfig struct {
	GrantDefaultExpiryDaysSinceLastUpload uint8 `mapstructure:"days_since_last_upload"`
	GrantDefaultExpiryTotalUploads        uint8 `mapstructure:"total_uploads"`
	GrantDefaultExpiryTotalDays           uint8 `mapstructure:"total_days"`
}

type ColorsConfig struct {
	Color       string     `mapstructure:"preset"`
	CustomColor [10]string `mapstructure:"custom_preset"`
}

type Config struct {
	OidcConfig        `mapstructure:"oidc"`
	DBConfig          `mapstructure:"db"`
	SMTPConfig        `mapstructure:"smtp"`
	FilesConfig       `mapstructure:"files"`
	ExpiryConfig      `mapstructure:"expiry"`
	GrantExpiryConfig `mapstructure:"grant_expiry"`
	LogConfig         `mapstructure:"log"`
	ColorsConfig      `mapstructure:"colors"`
	Otel              `mapstructure:"otel"`

	DevMode bool `mapstructure:"dev_mode"`

	Host           string   `mapstructure:"host"`
	Port           uint16   `mapstructure:"port"`
	RootPath       string   `mapstructure:"root_path"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

func setConfigSearchStrategy(viper *viper.Viper) {
	viper.SetConfigName("frans")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("frans")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func NewConfig() (Config, error) {
	var config Config
	fransConf := viper.New()

	fransConf.SetDefault("dev_mode", false)
	fransConf.SetDefault("root_path", "")

	fransConf.SetDefault("host", "127.0.0.1")
	fransConf.SetDefault("port", 8080)
	fransConf.SetDefault("trusted_proxies", []string{})

	fransConf.SetDefault("files.dir", "files")
	fransConf.SetDefault("files.max_per_upload", 20)
	fransConf.SetDefault("files.max_size", 2_000_000_000) // 2GB

	fransConf.SetDefault("expiry.days_since_last_download", 7)
	fransConf.SetDefault("expiry.total_downloads", 10)
	fransConf.SetDefault("expiry.total_days", 30)

	fransConf.SetDefault("grant_expiry.days_since_last_upload", 7)
	fransConf.SetDefault("grant_expiry.total_uploads", 10)
	fransConf.SetDefault("grant_expiry.total_days", 30)

	setDBConfigDefaults(fransConf)

	fransConf.SetDefault("oidc.issuer", "")
	fransConf.SetDefault("oidc.client_id", "")
	fransConf.SetDefault("oidc.admin_group", "admin")

	setLogConfigDefaults(fransConf)

	fransConf.SetDefault("smtp.server", "")
	fransConf.SetDefault("smtp.port", 25)
	fransConf.SetDefault("smtp.from", "")
	fransConf.SetDefault("smtp.username", nil)
	fransConf.SetDefault("smtp.password", nil)

	fransConf.SetDefault("colors.preset", "blue")
	fransConf.SetDefault("colors.custom_preset", [10]string{
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
		"#000000",
	})

	setOtelConfigDefaults(fransConf)

	setConfigSearchStrategy(fransConf)

	if err := fransConf.ReadInConfig(); err != nil {
		slog.Warn("No config file found, falling back to environment variables.")
	}

	if err := fransConf.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}

func (c *Config) GetBaseURL(request *http.Request) string {
	proto := "http"
	if request.TLS != nil {
		proto = "https"
	}
	host := request.Host
	patchedRootPath := c.RootPath
	if len(patchedRootPath) == 0 {
		patchedRootPath = "/"
	}
	return fmt.Sprintf("%s://%s%s", proto, host, patchedRootPath)
}
