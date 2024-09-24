set positional-arguments

[private]
help:
    just -l

# Runs all checks
check: generate fmt lint test-build test check-vulns

deps:
    go install go.uber.org/mock/mockgen@latest
    go install github.com/bufbuild/buf/cmd/buf@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install mvdan.cc/gofumpt@latest
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install gorm.io/gen/tools/gentool@latest
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Generate mocks
generate: build-proto generate-sql
    find . -type d -name "mocks" | xargs rm -rf
    go generate ./...

build-proto:
    #!/usr/bin/env bash
    cd pkg/application/proto
    buf mod update
    buf generate

# Lints the project
lint:
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
    go test -run ^$ ./... 1>/dev/null

# Runs all tests
test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p --race --cover

fmt:
    gofumpt -w -l .


bump:
    go get -u -v -t ./...
    go mod tidy

mkmigration MIGRATION_NAME:
    migrate create -ext sql -dir pkg/persistence/sqlite/internal/migrations -seq {{ MIGRATION_NAME }}

generate-sql:
    migrate -path pkg/persistence/sqlite/internal/migrations -database sqlite3:///tmp/db.db up
    gentool -db sqlite -dsn "file:///tmp/db.db?_fk=1&mode=ro" -outPath pkg/persistence/sqlite/internal/sqlgen -onlyModel
    rm /tmp/db.db

check-vulns:
    govulncheck ./...
    trivy fs . --exit-code 1 --ignore-unfixed --quiet
