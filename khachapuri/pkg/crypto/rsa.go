package crypto

import (
	native_crypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func RsaEncryptBytes(publicKey *RsaPublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		publicKey,
		data,
		nil)
}

func RsaDecryptBytes(privateKey *RsaPrivateKey, data []byte) ([]byte, error) {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA512 to hash the input.
	return privateKey.Decrypt(nil, data, &rsa.OAEPOptions{Hash: native_crypto.SHA512})
}

// LoadPublicKeyFromString loads a public RSA key from a string
func RsaLoadPublicKeyFromString(publicKeyStr string) (*RsaPublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Try to parse as PKIX first
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		if rsaPub, ok := pub.(*RsaPublicKey); ok {
			return rsaPub, nil
		}
		return nil, fmt.Errorf("not an RSA public key")
	}

	// Try to parse as PKCS1
	rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return rsaPub, nil
	}

	return nil, fmt.Errorf("failed to parse public key: %v", err)
}

func RsaLoadPrivateKeyFromString(privateKeyStr string) (*RsaPrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
