package aes

import (
	"crypto/aes"
	"encoding/hex"
)

var KEY = "lkajsfKLJFDFkljfasdfk213jldg453jAFd"

func EncryptAES(plaintext string) string {
	key := []byte(KEY)
	// create cipher
	c, _ := aes.NewCipher(key)

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}
