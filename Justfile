
default:
	just -l

# Installs binary dependencies of the project
deps:
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
