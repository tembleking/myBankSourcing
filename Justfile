
default:
    just -l

# Generate mocks
generate:
    go install github.com/golang/mock/mockgen@latest
    go generate ./...

lint:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    golangci-lint run --timeout 1h


ginkgo-bootstrap NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo bootstrap
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

ginkgo-generate NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

test-build:
    go build ./...

test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p
