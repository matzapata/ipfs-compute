package domain

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
)

type IRegistryService interface {
	ResolveDomain(domain string) (*contracts.Resolver, error)
	ResolveServer(domain string) (string, error)
	RegisterDomain(privateKey *ecdsa.PrivateKey, domain string, resolverAddress common.Address) (string, error)
	HashDomain(input string) [32]byte
}
