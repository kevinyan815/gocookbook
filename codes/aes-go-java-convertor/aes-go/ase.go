package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func Encrypt(key string, val string) (string, error) {
	origData := []byte(val)
	crypted, err := encrypt(key, origData)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(crypted), nil
}

func encrypt(key string, origData []byte) ([]byte, error) {
	if len(origData) <= 0 {
		return nil, errors.New("crypted len is zero")
	}
	keyBytes := GetKeyBytes(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func GetKeyBytes(key string) []byte {
	keyBytes := []byte(key)
	switch l := len(keyBytes); {
	case l < 16:
		keyBytes = append(keyBytes, make([]byte, 16-l)...)
	case l > 16:
		keyBytes = keyBytes[:16]
	}
	return keyBytes
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func main() {
	encryptKey := "3wjyxqDPNyrd4QrhxTycRMU4dFN2lCm4"
	sign, _ := Encrypt(encryptKey, "37b63ec62ebf8b2e790b8d9829da2ec26f1fad67")
	fmt.Println(sign)
}
