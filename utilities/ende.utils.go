package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

func Encrypt(plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plainText, nil), nil
}
func Decrypted(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	return gcm.Open(nil, nonce, cipherText, nil)
}
