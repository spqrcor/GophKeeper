package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// Encrypt шифрование сообщения открытым ключом
func Encrypt(input []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, input, nil)
	return ciphertext, nil
}

// Decrypt дешифрование сообщения закрытым ключом
func Decrypt(input []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := input[:nonceSize], input[nonceSize:]

	data, err := gcm.Open(nil, []byte(nonce), ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
