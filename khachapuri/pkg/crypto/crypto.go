package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/ethereum/go-ethereum/common"
)

type RsaPrivateKey = rsa.PrivateKey

type RsaPublicKey = rsa.PublicKey

type EcdsaPrivateKey = ecdsa.PrivateKey

type EcdsaPublicKey = ecdsa.PublicKey

type EcdsaAddress = common.Address

type EcdsaSignature struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}
