package lark

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// Decrypter ...
type Decrypter interface {
	DecryptString(s string) (res string, err error)
}

// EncryptEntry ...
type EncryptEntry struct {
	EncryptedBody string `json:"encrypt"`
}

type crypt struct {
	key []byte
}

// NewCrypto ...
func NewCrypto(keys string) Decrypter {
	k := sha256.Sum256([]byte(keys))
	return &crypt{key: k[:]}
}

// DecryptString ...
func (cd *crypt) DecryptString(s string) (res string, err error) {

	logger().Infow("decrypt string", "s", s)
	cipherText, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(cipherText)%aes.BlockSize != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}

	block, err := aes.NewCipher(cd.key)
	if err != nil {
		return
	}
	decryptedText := make([]byte, len(cipherText))
	decryptStream := cipher.NewCBCDecrypter(block, iv)
	decryptStream.CryptBlocks(decryptedText, cipherText)

	res = string(unpad(decryptedText, aes.BlockSize))
	return
}

func unpad(s []byte, size int) []byte {
	n := len(s)
	if n%size != 0 {
		return s
	}
	return s[:n-int(s[n-1])]
}
