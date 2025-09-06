package db

import (
	"fmt"

	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
)

func NewDBClient(dbConfig config.DBConfig) (*ent.Client, error) {
	var connString string
	switch dbConfig.DBType {
	case "postgres":
		if dbConfig.DBPort == 0 {
			dbConfig.DBPort = 5432
		}
		connString = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			dbConfig.DBHost,
			dbConfig.DBPort,
			dbConfig.DBUser,
			dbConfig.DBName,
			dbConfig.DBPassword,
		)
	case "mysql":
		if dbConfig.DBPort == 0 {
			dbConfig.DBPort = 3306
		}
		connString = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=True",
			dbConfig.DBUser,
			dbConfig.DBPassword,
			dbConfig.DBHost,
			dbConfig.DBPort,
			dbConfig.DBName,
		)
	case "sqlite3":
		if dbConfig.DBHost == "localhost" {
			dbConfig.DBHost = "frans.db"
		}
		connString = fmt.Sprintf("file:%s?cache=shared&_fk=1", dbConfig.DBHost)
	default:
		return nil, fmt.Errorf("database type %s is not supported by frans", dbConfig.DBType)
	}

	client, err := ent.Open(dbConfig.DBType, connString)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return client, nil
}
