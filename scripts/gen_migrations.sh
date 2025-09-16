#!/bin/bash
set -e
if [ $# -lt 1 ]; then
    echo "Usage: $0 <arg>" >&2
    exit 1
fi

MIGRATION_NAME=$1
echo "Generation migrations with name \"$MIGRATION_NAME\""

echo "Generating sqlite3..."
atlas migrate diff "$MIGRATION_NAME" --dir "file://internal/db/migrations/sqlite3" --to "ent://internal/ent/schema" --dev-url "sqlite://frans_migration?mode=memory&_fk=1"

echo "Generating postgres..."
atlas migrate diff "$MIGRATION_NAME" --dir "file://internal/db/migrations/postgres" --to "ent://internal/ent/schema" --dev-url "docker://postgres/17/test?search_path=public"

echo "Generating mysql..."
atlas migrate diff "$MIGRATION_NAME" --dir "file://internal/db/migrations/mysql" --to "ent://internal/ent/schema" --dev-url "docker://mysql/8/ent"

echo "Done!"