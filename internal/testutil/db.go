package testutil

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/db"
	_ "github.com/jvllmr/frans/internal/db/sqlite3"
	"github.com/jvllmr/frans/internal/ent"
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
