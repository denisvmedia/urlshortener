package shortener_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShortener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shortener Suite")
}
