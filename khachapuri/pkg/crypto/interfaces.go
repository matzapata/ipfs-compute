package crypto_service

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/ethereum/go-ethereum/common"
)

type ICryptoRsaService interface {
	EncryptBytes(publicKey *rsa.PublicKey, data []byte) ([]byte, error)
	DecryptBytes(privateKey *rsa.PrivateKey, data []byte) ([]byte, error)
	LoadPublicKeyFromString(publicKeyStr string) (*rsa.PublicKey, error)
	LoadPrivateKeyFromString(privateKeyStr string) (*rsa.PrivateKey, error)
}

type ICryptoEcdsaService interface {
	SignMessage(data []byte, hexkey string) (*Signature, error)
	LoadPrivateKeyFromString(hexkey string) (*ecdsa.PrivateKey, error)
	PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error)
}
