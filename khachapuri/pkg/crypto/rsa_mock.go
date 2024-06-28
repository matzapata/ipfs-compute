package crypto_service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MockCryptoRsaService struct {
	mock.Mock
}

func NewMockCryptoRsaService() *MockCryptoRsaService {
	return &MockCryptoRsaService{}
}

func (*MockCryptoRsaService) EncryptBytes(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	return []byte("encrypted"), nil
}

func (*MockCryptoRsaService) DecryptBytes(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	return []byte("decrypted"), nil
}

// LoadPublicKeyFromString loads a public RSA key from a string
func (*MockCryptoRsaService) LoadPublicKeyFromString(publicKeyStr string) (*rsa.PublicKey, error) {
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

func (*MockCryptoRsaService) LoadPrivateKeyFromString(privateKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
