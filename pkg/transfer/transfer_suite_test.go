package transfer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTransfer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transfer Suite")
}
