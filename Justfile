set positional-arguments

[private]
help:
    just -l

# Runs all checks
check: generate lint test-build test

# Generate mocks
generate:
    go install github.com/golang/mock/mockgen@latest
    find . -type d -name "mocks" | xargs rm -rf
    go generate ./...

# Lints the project
lint:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    golangci-lint run --timeout 1h

# Generates a ginkgo test in the current folder and bootstraps the suite if it doesn't exist
mktest NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo bootstrap 2>/dev/null || true
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

mkcommand *args='':
    go run github.com/spf13/cobra-cli@latest add "$@"

# Tests if the project builds correctly
test-build:
    go build ./...

# Runs all tests
test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p

# Launches the database in a docker container for development
db-launch: db-kill
    docker run --rm -d --pull always --name surrealdb -p 8000:8000 surrealdb/surrealdb:latest start --log trace --user root --pass root memory

db-logs:
    docker logs -f surrealdb

db-kill:
    docker rm -f surrealdb
