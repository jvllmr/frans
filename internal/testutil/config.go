package testutil

import (
	"log"
	"os"
	"strconv"

	"codeberg.org/jvllmr/frans/internal/config"
)

func SetupTestConfig() config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Could not parse config: %v", err)
	}
	cfg.FilesDir = "testfiles"

	smtpServer := os.Getenv("FRANS_SMTP_SERVER")
	if smtpServer == "" {
		cfg.SMTPServer = "127.0.0.1"
	} else {
		cfg.SMTPServer = smtpServer
	}

	smtpPort := os.Getenv("FRANS_SMTP_PORT")
	if smtpPort == "" {
		cfg.SMTPPort = 2525
	} else {
		num, err := strconv.Atoi(smtpPort)
		if err != nil {
			panic(err)
		}
		cfg.SMTPPort = num
	}

	cfg.SMTPFrom = "test_sender@vllmr.dev"
	return cfg
}
