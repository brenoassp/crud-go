package cryptography

type Provider interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}
