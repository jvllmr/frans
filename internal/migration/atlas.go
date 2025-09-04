package migration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"ariga.io/atlas/sql/migrate"
)

type scannable interface {
	Scan(dest ...any) error
}

type EntRevisionsReadWriter struct {
	db     *sql.DB
	dbType string
}

func (e *EntRevisionsReadWriter) formatSQLParam(count int) string {
	if e.dbType == "postgres" {
		return fmt.Sprintf("$%d", count)
	}

	return "?"
}

// DeleteRevision implements migrate.RevisionReadWriter.
func (e *EntRevisionsReadWriter) DeleteRevision(ctx context.Context, v string) error {
	_, err := e.db.QueryContext(
		ctx,
		fmt.Sprintf("DELETE FROM atlas_schema_revisions WHERE version = %s", e.formatSQLParam(1)),
		v,
	)
	return err
}

// Ident implements migrate.RevisionReadWriter.
func (e *EntRevisionsReadWriter) Ident() *migrate.TableIdent {
	return &migrate.TableIdent{
		Name:   "atlas_schema_revisions",
		Schema: "",
	}
}

func (*EntRevisionsReadWriter) scanRow(row scannable) (*migrate.Revision, error) {
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
func (e *EntRevisionsReadWriter) ReadRevision(
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
	return e.scanRow(row)
}

// ReadRevisions implements migrate.RevisionReadWriter.
func (e *EntRevisionsReadWriter) ReadRevisions(context.Context) ([]*migrate.Revision, error) {
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
func (e *EntRevisionsReadWriter) WriteRevision(ctx context.Context, rev *migrate.Revision) error {
	_, err := e.db.ExecContext(ctx, fmt.Sprintf(`
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
			) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)`,
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
		rev.PartialHashes,
		rev.OperatorVersion,
	)
	return err
}

var _ migrate.RevisionReadWriter = (*EntRevisionsReadWriter)(nil)
