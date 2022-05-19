package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/overlorddamygod/go-auth/configs"
)

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text string) (string, error) {
	mySecret := configs.GetConfig().TokenSecret1
	bytes := configs.GetConfig().TokenSecret2

	block, err := aes.NewCipher(mySecret)
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text string) (string, error) {
	mySecret := configs.GetConfig().TokenSecret1
	bytes := configs.GetConfig().TokenSecret2

	block, err := aes.NewCipher(mySecret)
	if err != nil {
		return "", err
	}
	cipherText, err := Decode(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
