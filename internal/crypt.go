package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
)

func Digest(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)

	return mac.Sum(nil)
}

func ValidateHMAC(message, digest, key []byte) bool {
	return hmac.Equal(digest, Digest(message, key))
}

func Crypt(contents, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	text := make([]byte, len(contents))

	// See: https://github.com/fugue/credstash/issues/75
	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(text, contents)

	return text
}
