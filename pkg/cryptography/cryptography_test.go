package cryptography

import(
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cryptography", func() {
	var message string
	BeforeEach(func() {
		message = "sample message"
	})

	When("Message is sent to the user using its public key", func() {
		It("Decrypt the message", func() {
			privateKey := generatePrivateKey()
			publicKey := privateKey.PublicKey
			ciphertext := EncryptMessage(message, publicKey)
			text := DecryptMessage(ciphertext, privateKey)
			Expect(text).To(Equal(message))
		})
	})
})