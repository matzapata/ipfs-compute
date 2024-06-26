package crypto_service

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signature struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

type CryptoEcdsaService struct {
}

func NewCryptoEcdsaService() *CryptoEcdsaService {
	return &CryptoEcdsaService{}
}

func (*CryptoEcdsaService) SignMessage(data []byte, hexkey string) (*Signature, error) {
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

func (*CryptoEcdsaService) LoadPrivateKeyFromString(hexkey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexkey)
}

func (*CryptoEcdsaService) PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
