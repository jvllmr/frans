package cmd

import (
	"log"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/db"
	"github.com/jvllmr/frans/internal/ent"
)

func getConfigAndDBClient() (config.Config, *ent.Client) {
	configValue := config.NewSafeConfig()
	db, err := db.NewDBClient(configValue.DBConfig)
	if err != nil {
		log.Fatalf("Could not create db client: %v", err)
	}
	return configValue, db
}
