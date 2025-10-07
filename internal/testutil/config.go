package testutil

import (
	"log"

	"github.com/jvllmr/frans/internal/config"
)

func SetupTestConfig() config.Config {
	configValue, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Could not parse config: %v", err)
	}
	configValue.FilesDir = "testfiles"
	configValue.SMTPServer = "127.0.0.1"
	configValue.SMTPPort = 2525
	configValue.SMTPFrom = "test_sender@vllmr.dev"
	return configValue
}
