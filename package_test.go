package testbasher

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTestBasher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "testbasher package")
}
