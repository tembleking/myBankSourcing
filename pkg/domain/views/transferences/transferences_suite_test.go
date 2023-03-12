package transferences_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTransferences(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transferences Suite")
}
