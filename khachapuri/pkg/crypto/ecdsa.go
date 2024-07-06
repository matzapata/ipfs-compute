package crypto

import (
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	geth_crypto "github.com/ethereum/go-ethereum/crypto"
)

func EcdsaSignMessage(data []byte, hexkey string) (*EcdsaSignature, error) {
	privateKey, err := geth_crypto.HexToECDSA(hexkey)
	if err != nil {
		log.Fatal(err)
	}

	hash := geth_crypto.Keccak256Hash(data)
	signature, err := geth_crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return &EcdsaSignature{
		Hash:      hash.Hex(),
		Signature: hexutil.Encode(signature),
		Address:   geth_crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
	}, nil
}

func EcdsaHexToPrivateKey(hexkey string) *EcdsaPrivateKey {
	pk, err := geth_crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(err)
	}

	return pk
}

func EcdsaPrivateKeyToAddress(privateKey *EcdsaPrivateKey) *EcdsaAddress {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*EcdsaPublicKey)
	if !ok {
		panic(errors.New("error casting public key to ECDSA"))
	}

	addr := geth_crypto.PubkeyToAddress(*publicKeyECDSA)
	return &addr
}

func EcdsaHexToAddress(hex string) *EcdsaAddress {
	addr := common.HexToAddress(hex)
	return &addr
}
