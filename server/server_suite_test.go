package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUrlshortener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UrlShortener Suite")
}
