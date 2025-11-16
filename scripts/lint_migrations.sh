set -e

atlas migrate lint --dir "file://internal/db/migrations/sqlite3" --dev-url "sqlite://frans_migration?mode=memory&_fk=1" --git-base main


atlas migrate lint --dir "file://internal/db/migrations/postgres" --dev-url "docker://postgres/17/test?search_path=public" --git-base main


atlas migrate lint --dir "file://internal/db/migrations/mysql" --dev-url "docker://mysql/8/ent" --git-base main
