package surrealdb_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSurrealdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Surrealdb Suite")
}
