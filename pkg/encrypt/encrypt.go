package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type EncryptTool interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

type encryptTool struct {
	secretKey string
	iv        string
}

func NewEncryptTool(secretKey, iv string) *encryptTool {
	return &encryptTool{
		secretKey: secretKey,
		iv:        iv,
	}
}

func (e *encryptTool) Encrypt(text string) (string, error) {
	key := e.secretKey
	iv := e.iv

	var plainTextBlock []byte
	length := len(text)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, text)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func (e *encryptTool) Decrypt(text string) (string, error) {
	key := e.secretKey
	iv := e.iv

	ciphertext, err := base64.StdEncoding.DecodeString(text)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return string(ciphertext), nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
