package migration

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/sqlite"
	"github.com/jvllmr/frans/internal/config"

	"github.com/jvllmr/frans/internal/util"
)

//go:embed migrations/*/*.sql
//go:embed migrations/*/*.sum
var migrationFiles embed.FS

func buildMigrationDsn(dbConfig config.DBConfig) (string, error) {
	var dsn string

	switch dbConfig.DBType {
	case "postgres":
		if dbConfig.DBPort == 0 {
			dbConfig.DBPort = 5432
		}
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?search_path=public&sslmode=disable",
			dbConfig.DBUser,
			dbConfig.DBPassword,
			dbConfig.DBHost,
			dbConfig.DBPort,
			dbConfig.DBName,
		)
	case "mysql":
		if dbConfig.DBPort == 0 {
			dbConfig.DBPort = 3306
		}
		dsn = fmt.Sprintf(
			"mysql://%s:%s@%s:%d/%s",
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
		dsn = fmt.Sprintf("file:%s?_fk=1", dbConfig.DBHost)
	default:
		return "", fmt.Errorf("database type %s is not supported by frans", dbConfig.DBType)
	}

	return dsn, nil
}

func getDriver(db *sql.DB, dbType string) (migrate.Driver, error) {
	var drv migrate.Driver
	var err error
	switch dbType {
	case "postgres":
		drv, err = postgres.Open(db)
		if err != nil {
			return nil, err
		}

	case "mysql":
		drv, err = mysql.Open(db)
		if err != nil {
			return nil, err
		}

	case "sqlite3":
		drv, err = sqlite.Open(db)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("database type %s is not supported by frans", dbType)
	}

	return drv, err
}

func Migrate() {
	dbConfig, err := config.GetDBConfig()
	if err != nil {
		log.Fatalf("Could not get database config: %v", err)
	}
	dsn, err := buildMigrationDsn(dbConfig)

	if err != nil {
		log.Fatalf("Could not build database dsn: %v", err)
	}

	db, err := sql.Open(dbConfig.DBType, dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	if err := util.UnpackFSToPath(migrationFiles, "."); err != nil {
		log.Fatalf("Could not unpack migration files: %v", err)
	}
	defer func() {
		if err := os.RemoveAll("migrations"); err != nil {
			log.Printf("Could not cleanup all migration files: %v", err)
		}
	}()

	drv, err := getDriver(db, dbConfig.DBType)
	if err != nil {
		log.Fatalf("Could not get database driver: %v", err)
	}

	dir, err := migrate.NewLocalDir(fmt.Sprintf("migrations/%s", dbConfig.DBType))

	if err != nil {
		log.Fatalf("Could not create new local dir: %v", err)
	}

	rrw := entRevisionsReadWriter{db: db, dbType: dbConfig.DBType}
	err = rrw.createTable()
	if err != nil {
		log.Fatalf("Could not create revisions table: %v", err)
	}

	executor, err := migrate.NewExecutor(drv, dir, &rrw, migrate.WithAllowDirty(true))
	if err != nil {
		log.Fatalf("Could not get migration executor: %v", err)
	}
	err = executor.ExecuteN(context.Background(), 0)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoPendingFiles) {
			log.Fatalf("Could not execute pending migrations: %v", err)
		}
	}

}
