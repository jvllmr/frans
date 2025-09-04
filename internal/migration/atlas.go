package migration

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"ariga.io/atlas/sql/migrate"
)

type scannable interface {
	Scan(dest ...any) error
}

type entRevisionsReadWriter struct {
	db     *sql.DB
	dbType string
}

func (e *entRevisionsReadWriter) formatSQLParam(count int) string {
	if e.dbType == "postgres" {
		return fmt.Sprintf("$%d", count)
	}

	return "?"
}

// DeleteRevision implements migrate.RevisionReadWriter.
func (e *entRevisionsReadWriter) DeleteRevision(ctx context.Context, v string) error {
	_, err := e.db.QueryContext(
		ctx,
		fmt.Sprintf("DELETE FROM atlas_schema_revisions WHERE version = %s", e.formatSQLParam(1)),
		v,
	)
	return err
}

// Ident implements migrate.RevisionReadWriter.
func (e *entRevisionsReadWriter) Ident() *migrate.TableIdent {
	return &migrate.TableIdent{
		Name:   "atlas_schema_revisions",
		Schema: "",
	}
}

func (*entRevisionsReadWriter) scanRow(row scannable) (*migrate.Revision, error) {
	var rev migrate.Revision
	var partialHashesHolder []byte
	err := row.Scan(
		&rev.Version,
		&rev.Description,
		&rev.Type,
		&rev.Applied,
		&rev.Total,
		&rev.ExecutedAt,
		&rev.ExecutionTime,
		&rev.Error,
		&rev.ErrorStmt,
		&rev.Hash,
		&partialHashesHolder,
		&rev.OperatorVersion,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(partialHashesHolder, &rev.PartialHashes)

	if err != nil {
		return nil, err
	}

	return &rev, nil
}

// ReadRevision implements migrate.RevisionReadWriter.
func (e *entRevisionsReadWriter) ReadRevision(
	ctx context.Context,
	v string,
) (*migrate.Revision, error) {

	row := e.db.QueryRowContext(context.Background(), fmt.Sprintf(`SELECT 
			version, 
			description, 
			type, 
			applied, 
			total, 
			executed_at, 
			execution_time, 
			error, 
			error_stmt, 
			hash, 
			partial_hashes, 
			operator_version
			FROM atlas_schema_revisions
			WHERE version = %s`,
		e.formatSQLParam(1),
	),
		v,
	)
	rev, err := e.scanRow(row)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, migrate.ErrRevisionNotExist
	}
	return rev, err
}

// ReadRevisions implements migrate.RevisionReadWriter.
func (e *entRevisionsReadWriter) ReadRevisions(context.Context) ([]*migrate.Revision, error) {
	rows, err := e.db.QueryContext(context.Background(), `SELECT
			version, 
			description, 
			type, 
			applied, 
			total, 
			executed_at, 
			execution_time, 
			error, 
			error_stmt, 
			hash, 
			partial_hashes, 
			operator_version
			FROM atlas_schema_revisions
			`,
	)
	if err != nil {
		return nil, err
	}

	revs := make([]*migrate.Revision, 0)
	for rows.Next() {
		rev, err := e.scanRow(rows)
		if err != nil {
			return nil, err
		}
		revs = append(revs, rev)
	}

	return revs, nil
}

// WriteRevision implements migrate.RevisionReadWriter.
func (e *entRevisionsReadWriter) WriteRevision(ctx context.Context, rev *migrate.Revision) error {
	encPartialHashes, err := json.Marshal(rev.PartialHashes)
	if err != nil {
		return err
	}

	_, err = e.db.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO atlas_schema_revisions (
			version, 
			description, 
			type, 
			applied, 
			total, 
			executed_at, 
			execution_time, 
			error, 
			error_stmt, 
			hash, 
			partial_hashes, 
			operator_version
			) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
			 ON CONFLICT(version) DO UPDATE
			SET description = %s, 
				type = %s, 
				applied = %s, 
				total = %s, 
				executed_at = %s, 
				execution_time = %s, 
				error = %s, 
				error_stmt = %s, 
				hash = %s, 
				partial_hashes = %s, 
				operator_version = %s`,
		e.formatSQLParam(1),
		e.formatSQLParam(2),
		e.formatSQLParam(3),
		e.formatSQLParam(4),
		e.formatSQLParam(5),
		e.formatSQLParam(6),
		e.formatSQLParam(7),
		e.formatSQLParam(8),
		e.formatSQLParam(9),
		e.formatSQLParam(10),
		e.formatSQLParam(11),
		e.formatSQLParam(12),
		e.formatSQLParam(13),
		e.formatSQLParam(14),
		e.formatSQLParam(15),
		e.formatSQLParam(16),
		e.formatSQLParam(17),
		e.formatSQLParam(18),
		e.formatSQLParam(19),
		e.formatSQLParam(20),
		e.formatSQLParam(21),
		e.formatSQLParam(22),
		e.formatSQLParam(23),
	),
		rev.Version,
		rev.Description,
		rev.Type,
		rev.Applied,
		rev.Total,
		rev.ExecutedAt,
		rev.ExecutionTime,
		rev.Error,
		rev.ErrorStmt,
		rev.Hash,
		encPartialHashes,
		rev.OperatorVersion,
		rev.Description,
		rev.Type,
		rev.Applied,
		rev.Total,
		rev.ExecutedAt,
		rev.ExecutionTime,
		rev.Error,
		rev.ErrorStmt,
		rev.Hash,
		encPartialHashes,
		rev.OperatorVersion,
	)
	return err
}

func (e *entRevisionsReadWriter) createTable() error {
	var err error
	switch e.dbType {
	case "postgres":
		_, err = e.db.Exec(`CREATE TABLE IF NOT EXISTS (
		version
	)`)
	case "mysql":
		_, err = e.db.Exec(`CREATE TABLE IF NOT EXISTS (
		version
	)`)

	case "sqlite3":
		_, err = e.db.Exec(`CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
		version text NOT NULL, 
		description text NOT NULL, 
		type integer NOT NULL DEFAULT (2), 
		applied integer NOT NULL DEFAULT (0),
		total integer NOT NULL DEFAULT (0), 
		executed_at datetime NOT NULL, 
		execution_time integer NOT NULL, 
		error text NULL, 
		error_stmt text NULL, 
		hash text NOT NULL, 
		partial_hashes json NULL, 
		operator_version text NOT NULL,
		PRIMARY KEY (version)
	)`)
	default:
		err = fmt.Errorf("database type %s is not supported by frans", e.dbType)
	}

	return err
}

var _ migrate.RevisionReadWriter = (*entRevisionsReadWriter)(nil)
