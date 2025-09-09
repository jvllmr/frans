package sqlite3

import (
	"database/sql"

	"modernc.org/sqlite"
)

// See https://github.com/ent/ent/issues/2460
func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}
