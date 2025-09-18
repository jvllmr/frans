package testutil

import "github.com/jvllmr/frans/internal/config"

func SetupTestConfig(configModifier func(configValue *config.Config) *config.Config) config.Config {
	if configModifier == nil {
		configModifier = func(configValue *config.Config) *config.Config { return configValue }
	}
	configValue := config.NewSafeConfig()
	configValue.FilesDir = "testfiles"
	return *configModifier(&configValue)
}
