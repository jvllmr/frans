#!/bin/bash

MIGRATION_NAME=$1

atlas migrate diff "$MIGRATION_NAME" --dir "file://internal/migration/migrations" --to "ent://internal/ent/schema" --dev-url "sqlite://frans?mode=memory&_fk=1"