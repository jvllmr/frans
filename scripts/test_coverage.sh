#!/bin/bash
set -e
go test -coverprofile=coverage.txt \
    "./internal/routes/api" \
    "./internal/tasks" \
    "./internal/routes/api/share" \
    -coverpkg=./internal/services,\
./internal/db,\
./internal/config,\
./internal/util,\
./internal/testutil,\
./internal/routes/api/types,\
./internal/mail,\
./internal/routes/api,\
./internal/tasks,\
./internal/routes/api/share