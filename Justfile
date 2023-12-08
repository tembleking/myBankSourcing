set positional-arguments

[private]
help:
    just -l

# Runs all checks
check: generate lint test-build test

# Generate mocks
generate: build-proto generate-sql
    go install github.com/golang/mock/mockgen@latest
    find . -type d -name "mocks" | xargs rm -rf
    go generate ./...

build-proto:
    #!/usr/bin/env bash
    go install github.com/bufbuild/buf/cmd/buf@latest
    cd pkg/application/proto
    buf mod update
    buf generate

# Lints the project
lint:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    golangci-lint run --timeout 1h

# Generates a ginkgo test in the current folder and bootstraps the suite if it doesn't exist
mktest NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo bootstrap 2>/dev/null || true
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

# Adds a new command to the project with the given name and parent command.
mkcommand NAME PARENT_COMMAND:
    cd ./cmd/clerk; go run github.com/spf13/cobra-cli@latest add {{ NAME }} --parent {{ PARENT_COMMAND }}

# Tests if the project builds correctly
test-build:
    go build ./...
    go test -run ^$ ./...

# Runs all tests
test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p --race --cover

mkmigration MIGRATION_NAME:
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    migrate create -ext sql -dir pkg/persistence/sqlite/internal/migrations -seq {{ MIGRATION_NAME }}

generate-sql:
    go install github.com/go-jet/jet/v2/cmd/jet@latest
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    migrate -path pkg/persistence/sqlite/internal/migrations -database sqlite3:///tmp/db.db up
    jet -source sqlite -dsn file:///tmp/db.db -path pkg/persistence/sqlite/internal/sqlgen
    rm /tmp/db.db
