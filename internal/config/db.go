package config

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/jvllmr/frans/internal/ent"
)

var DBClient *ent.Client

func InitDB(configValue Config) {
	var connString string
	switch configValue.DBType {
	case "postgres":
		if configValue.DBPort == 0 {
			configValue.DBPort = 5432
		}
		connString = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			configValue.DBHost,
			configValue.DBPort,
			configValue.DBUser,
			configValue.DBName,
			configValue.DBPassword,
		)
	case "mysql":
		if configValue.DBPort == 0 {
			configValue.DBPort = 3306
		}
		connString = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=True",
			configValue.DBUser,
			configValue.DBPassword,
			configValue.DBHost,
			configValue.DBPort,
			configValue.DBName,
		)
	case "sqlite3":
		if configValue.DBHost == "localhost" {
			configValue.DBHost = "frans.db"
		}
		connString = fmt.Sprintf("file:%s?cache=shared&_fk=1", configValue.DBHost)
	default:
		slog.Error(fmt.Sprintf("Database type %s is not supported by frans", configValue.DBType))
		os.Exit(1)
	}

	client, err := ent.Open(configValue.DBType, connString)
	if err != nil {
		slog.Error("failed opening connection to database", "err", err)
		os.Exit(1)
	}
	DBClient = client
}
