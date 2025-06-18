package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jvllmr/frans/pkg/ent"
)

var DBClient *ent.Client

func InitDB(configValue Config) {
	var connString string
	if configValue.DBType == "postgres" {
		if configValue.DBPort == 0 {
			configValue.DBPort = 5432
		}
		connString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", configValue.DBHost, configValue.DBPort, configValue.DBUser, configValue.DBName, configValue.DBPassword)
	} else if configValue.DBType == "mysql" {
		if configValue.DBPort == 0 {
			configValue.DBPort = 3306
		}
		connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", configValue.DBUser, configValue.DBPassword, configValue.DBHost, configValue.DBPort, configValue.DBName)
	} else if configValue.DBType == "sqlite3" {
		if configValue.DBHost == "localhost" {
			configValue.DBHost = "frans.db"
		}
		connString = fmt.Sprintf("file:%s?cache=shared&_fk=1", configValue.DBHost)
	} else {
		log.Fatalf("Error: Database type %s is not supported by frans", configValue.DBType)
	}

	client, err := ent.Open(configValue.DBType, connString)
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	if configValue.DevMode {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}
	DBClient = client
}
