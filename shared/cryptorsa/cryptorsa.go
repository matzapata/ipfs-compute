package cryptorsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func EncryptBytes(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		publicKey,
		data,
		nil)
}

// LoadPublicKeyFromString loads a public RSA key from a string
func LoadPublicKeyFromString(publicKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Try to parse as PKIX first
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
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

// LoadPublicKeyFromFile loads a public RSA key from a file
func LoadPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return LoadPublicKeyFromString(string(data))
}
