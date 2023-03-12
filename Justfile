
default:
    just -l

# Generate mocks
generate:
    go install github.com/golang/mock/mockgen@latest
    go generate ./...

ginkgo-bootstrap NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo bootstrap
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

ginkgo-generate NAME="":
    cd {{invocation_directory()}}; go run github.com/onsi/ginkgo/v2/ginkgo generate {{ NAME }}

test-build:
    go build ./...

test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p