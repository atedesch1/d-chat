package cryptography_test

import(
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/decentralized-chat/pkg/cryptography"
)

var _ = Describe("Cryptography", func() {
	var message string
	BeforeEach(func() {
		message = "sample message"
	})

	When("Message is sent to the user using its public key", func() {
		It("Decrypt the message", func() {
			privateKey := cryptography.generatePrivateKey()
			publicKey := privateKey.PublicKey
			ciphertext := cryptography.EncryptMessage(message, publicKey)
			text := cryptography.DecryptMessage(ciphertext, privateKey)
			Expect(text).To(Equal(message))
		})
	})
})