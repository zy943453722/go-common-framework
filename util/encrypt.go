package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func GetMd5(str string, salts ...string) string {
	h := md5.New()
	h.Write([]byte(str))
	if len(salts) == 0 {
		return hex.EncodeToString(h.Sum(nil))
	} else {
		for _, salt := range salts {
			h.Write([]byte(salt))
		}
		return hex.EncodeToString(h.Sum(nil))
	}
}

func Base64Encrypt(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}

func Base64Decrypt(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func GetSha512(str string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(str)))
}

//AesCBCEncrypt AES的CBC加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	//block大小和初始向量大小保持一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

//AesCBCDecrypt AES的CBC解密
func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData, err = PKCS7UnPadding(encryptData)
	if err != nil {
		return nil, err
	}
	return encryptData, nil
}

//PKCS7Padding 使用PKCS7进行填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

//PKCS7UnPadding 解PKCS7填充
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if length < unPadding {
		return nil, errors.New("ciphertext not valid")
	}
	return origData[:(length - unPadding)], nil
}
