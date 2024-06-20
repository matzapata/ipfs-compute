package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signature struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func SignMessage(data []byte, hexkey string) (*Signature, error) {
	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		log.Fatal(err)
	}

	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return &Signature{
		Hash:      hash.Hex(),
		Signature: hexutil.Encode(signature),
		Address:   crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
	}, nil
}

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
