package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

func generatePrivateKey() rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return *privateKey
}

func EncryptMessage(message string, publicKey rsa.PublicKey) string {
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &publicKey, []byte(message), nil)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptMessage(cipherString string, privateKey rsa.PrivateKey) string {
	ciphertext, _ := base64.StdEncoding.DecodeString(cipherString)
	rng := rand.Reader
	message, err := rsa.DecryptOAEP(sha256.New(), rng, &privateKey, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return string(message)
}