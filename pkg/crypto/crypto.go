package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

type crypto struct {
	key [32]byte
}

func NewCrypto(key string) *crypto {
	return &crypto{
		key: sha256.Sum256([]byte(key)),
	}
}

func (c *crypto) Encode(password string) (string, error) {
	aesblock, err := aes.NewCipher(c.key[:])
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce := c.key[len(c.key)-aesgcm.NonceSize():]

	encode := aesgcm.Seal(nil, nonce, []byte(password), c.key[:])

	return hex.EncodeToString(encode), nil
}
