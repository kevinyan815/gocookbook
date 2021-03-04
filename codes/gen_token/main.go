package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	aesKEY = "justTest!@#$test" // must be 16 characters
	md5Len = 4  //MD5 part size in access token
	aesLen = 16 //bytes after being encrypted by aes
)

func main() {
	token := genAccessToken(1234567)
	println(token)

	println(ParseAccessToken(token))
}

// generate access token
// + means concatenate
// Token = Hex(MD5(uid + time)[0:4] + AES(uid + time))
// a string with length of 40
func genAccessToken(uid int64) string {
	byteInfo := make([]byte, 12)
	binary.BigEndian.PutUint64(byteInfo, uint64(uid))
	binary.BigEndian.PutUint32(byteInfo[8:], uint32(time.Now().Unix()))
	encodeByte, err := AesEncrypt(byteInfo, []byte(aesKEY)) //returned slice's size is 16
	if err != nil {
		fmt.Println("desc: ", "genAccessToken-AesEncrypt", "error: ", err, "byteInfo: ", byteInfo)
	}
	md5Byte := md5.Sum(byteInfo)
	data := append(md5Byte[0:md5Len], encodeByte...)

	return hex.EncodeToString(data)
}

// parse uid from access token
func ParseAccessToken(accessToken string) (uid uint64) {
	if len(accessToken) != 2*(md5Len+aesLen) {
		fmt.Println("log_desc", "len(accessToken)", "length", len(accessToken))
		return
	}
	encodeStr := accessToken[md5Len*2:]
	data, e := hex.DecodeString(encodeStr)
	if e != nil {
		fmt.Println("log_desc", "ParseAccessToken-DecodeString)", "error", e, "accessToken", encodeStr)
		return
	}
	decodeByte, e := AesDecrypt(data, []byte(aesKEY)) //忽略错误
	if e != nil {
		fmt.Println("log_desc", "ParseAccessToken-AesDecrypt)", "error", e, "data", data)
	}
	uid = binary.BigEndian.Uint64(decodeByte)

	return uid
}

//key is 16 bytes
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)

	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)

	unPadding := int(origData[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return origData[:(length - unPadding)]
}
