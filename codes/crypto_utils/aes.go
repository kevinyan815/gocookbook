package crypto_utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// AES KEY either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
// var aseKey = "b3c65d06a1cf4dbda57af0af5c63e85f"

func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	res, err := base64.URLEncoding.DecodeString(data)
	fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)
}

func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

/*
//  keyStr 密钥
//  value  消息内容
*/
func HMACSHA1(value, keyStr string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}

// key is 32 bytes
func AesEcbPkcs5Decrypt(crypted, key []byte) (origData []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	blockMode := NewECBDecrypter(block)
	origData = make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return
}

func AesEcbPkcs5Encrypt(src string, key []byte) (crypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted = make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func AES256CBCEncrypt(plantText string, key string, iv []byte) (value string, err error) {
	plaintext := PKCS7Padding([]byte(plantText))
	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AES256CBCDecrypt(plantText string, key string, iv string) (value string, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return
	}
	var iiv []byte
	if iiv, err = base64.StdEncoding.DecodeString(iv); err != nil {
		return
	}

	var cipherText []byte
	if cipherText, err = base64.StdEncoding.DecodeString(plantText); err != nil {
		return
	}

	mode := cipher.NewCBCDecrypter(block, iiv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText = PKCS7UnPadding(cipherText)
	return string(cipherText), nil
}

//key is 16 bytes
func AesCbcPkcs5Encrypt(origData, key []byte) ([]byte, error) {
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

func AesCbcPkcs5Decrypt(crypted, key []byte) ([]byte, error) {
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

func getHmacCode(s string, key []byte) string {
	h := hmac.New(sha256.New, key)
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Crypt struct {
	Iv    string `json:"iv"`
	Value string `json:"value"`
	Mac   string `json:"mac"`
}

func OpenSslAesEncrypt(value string, key string) (payload string, err error) {
	iv := make([]byte, 16)
	_, err = rand.Read(iv)

	var encode string
	if encode, err = AES256CBCEncrypt(value, key, iv); err != nil {
		return
	}

	ivv := base64.StdEncoding.EncodeToString(iv)

	crypt := Crypt{Iv: ivv, Value: encode, Mac: getHmacCode(ivv+encode, []byte(key))}

	var bs []byte
	if bs, err = json.Marshal(crypt); err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

func OpenSslAesDecrypt(payload string, key string) (value string, err error) {
	var bs []byte
	if bs, err = base64.StdEncoding.DecodeString(payload); err != nil {
		return
	}
	crypt := Crypt{}
	if err = json.Unmarshal(bs, &crypt); err != nil {
		return
	}

	return AES256CBCDecrypt(crypt.Value, key, crypt.Iv)
}


// Demo Application
func aesEncryptDemo(content string) string {
	aesKey := "b3c65d06a1cf4dbda57af0af5c63e85f"
	aesKeyByte, _ := hex.DecodeString(aesKey)
	encrypt, _ := AesEcbPskcs5Encrypt(content, aesKeyByte)
	return hex.EncodeToString(encrypt)
}

func aesDecryptDemo(crypt string) (reply string, err error) {
	if crypt == "" {
		return
	}
	aesKey := "b3c65d06a1cf4dbda57af0af5c63e85f"
	aesKeyByte, _ := hex.DecodeString(aesKey)
	cryptByte, _ := hex.DecodeString(crypt)
	decryptByte, err := AesEcbPskcs5Decrypt(cryptByte, aesKeyByte)
	return string(decryptByte), err
}
