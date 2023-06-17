package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func EncryptAES(src string, key string, vector string) (string, error) {
	block, err := aes.NewCipher([]byte(key));
	if err != nil {
		return "", err;
	}

	if src == "" {
		return "", errors.New("Plain content empty");
	}

	ecb := cipher.NewCBCEncrypter(block, []byte(vector));
	content := []byte(src);
	content = PKCS5Padding(content, block.BlockSize());
	crypt := make([]byte, len(content));
	ecb.CryptBlocks(crypt, content);
	result := Base64EncodeToString(crypt);

	return result, nil;
}

func DecryptAES(src string, key string, vector string) ([]byte, error) {
	crypt, err := Base64DecodeString(src);
	if err != nil {
		return nil, err;
	}

	block, err := aes.NewCipher([]byte(key));
	if err != nil {
		return nil, err;
	}

	if len(crypt) == 0 {
		return nil, errors.New("Plain content empty");
	}

	ecb := cipher.NewCBCDecrypter(block, []byte(vector));
	decrypted := make([]byte, len(crypt));
	ecb.CryptBlocks(decrypted, crypt);
	return PKCS5Trimming(decrypted), err;
}

func Base64EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src);
}

func Base64DecodeString(src string) ([]byte, error) {
	decryptedText, err := base64.StdEncoding.DecodeString(src);
	if err != nil {
		return nil, err;
	}
	return decryptedText, nil;
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
	
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1];
	return encrypt[:len(encrypt)-int(padding)]
}