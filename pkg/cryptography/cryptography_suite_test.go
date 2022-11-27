package cryptography

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCryptography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cryptography Suite")
}
