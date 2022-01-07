package testbasher

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestBasher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "testbasher package")
}
