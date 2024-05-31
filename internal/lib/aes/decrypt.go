package aes

import (
	"crypto/aes"
	"encoding/hex"
)

func DecryptAES(ct string) string {
	key := []byte(KEY)
	ciphertext, _ := hex.DecodeString(ct)

	c, _ := aes.NewCipher(key)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s
}
