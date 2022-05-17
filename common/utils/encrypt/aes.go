package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesEncrypt(encodeStr string, aesKey string, aesIv string) (string, error) {
	return AesEncryptKeyIv(encodeStr, aesKey, aesIv)
}

func AesEncryptKeyIv(encodeStr, key, iv string) (string, error) {
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = pKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	//填充
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipherText, padText...)
}

func AesDecrypt(decodeStr string, aesKey string, aesIv string) ([]byte, error) {
	return AesDecryptKeyIv(decodeStr, aesKey, aesIv)
}

func CiphertextLength(n int) int {
	aesLen := aes.BlockSize * (n/aes.BlockSize + 1)
	size := aesLen / 3 * 4
	if aesLen%3 != 0 {
		size += 4
	}
	return size
}

func AesDecryptKeyIv(decodeStr, key, iv string) ([]byte, error) {
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
