package testutil

import (
	"os"
	"testing"

	"codeberg.org/jvllmr/frans/internal/config"
	"codeberg.org/jvllmr/frans/internal/db"
	_ "codeberg.org/jvllmr/frans/internal/db/sqlite3"
	"codeberg.org/jvllmr/frans/internal/ent"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDBClient(t *testing.T) *ent.Client {

	dbConfig := config.DBConfig{
		DBType: "sqlite3",
		DBHost: "test.db",
	}
	db.Migrate(dbConfig)
	db, err := db.NewDBClient(dbConfig)
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		db.Close()
		if err := os.Remove("test.db"); err != nil {
			panic(err)
		}
	})
	return db
}
