package domain

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

type IRegistryService interface {
	ResolveDomain(domain string) (cid string, err error)
	RegisterDomain(privateKey *ecdsa.PrivateKey, domain string, resolverAddress common.Address) (string, error)
	HashDomain(input string) [32]byte
}
