package testutil

import "github.com/jvllmr/frans/internal/config"

func SetupTestConfig() config.Config {
	configValue := config.NewSafeConfig()
	configValue.FilesDir = "testfiles"
	configValue.SMTPServer = "127.0.0.1"
	configValue.SMTPPort = 2525
	configValue.SMTPFrom = "test_sender@vllmr.dev"
	return configValue
}
