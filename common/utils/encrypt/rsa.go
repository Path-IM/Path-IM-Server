package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RSAKeyGenToKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {

		return nil, err

	}

	return privateKey, nil
}

func PubKeyToBytes(key *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(key)
}

func BytesToPubKey(key []byte) (interface{}, error) {
	return x509.ParsePKCS1PublicKey(key)
}

func RSAKeyGenToBuf(bits int) ([]byte, []byte, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {

		return nil, nil, err

	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	block := &pem.Block{

		Type: "RSA Private key",

		Bytes: derStream,
	}

	priBuf := &bytes.Buffer{}

	err = pem.Encode(priBuf, block)

	if err != nil {

		return nil, nil, err

	}

	publicKey := &privateKey.PublicKey

	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)

	block = &pem.Block{

		Type: "RSA Public key",

		Bytes: derPkix,
	}

	if err != nil {

		return nil, nil, err

	}

	pubBuf := &bytes.Buffer{}

	err = pem.Encode(pubBuf, block)

	if err != nil {

		return nil, nil, err

	}

	return priBuf.Bytes(), pubBuf.Bytes(), nil
}

func RsaPrivateKeyDecrypt(cipher []byte, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pri, cipher)
}
