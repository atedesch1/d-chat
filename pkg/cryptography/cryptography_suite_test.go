package cryptography_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCryptography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cryptography Suite")
}
