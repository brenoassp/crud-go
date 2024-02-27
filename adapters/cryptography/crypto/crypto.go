package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Client struct {
	key string
}

func NewClient(key string) Client {
	return Client{key: key}
}

func (c Client) Encrypt(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nonce, nonce, plainText, nil)

	return encryptedData, nil
}

func (c Client) Decrypt(cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, err
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
