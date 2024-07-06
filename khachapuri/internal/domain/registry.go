package domain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type ProviderDomainData struct {
	EcdsaAddress   string
	RsaPublicKey   string
	ServerEndpoint string
}

type IRegistryService interface {
	ResolveDomain(domainName string) (*ProviderDomainData, error)
	RegisterDomain(privateKey *crypto.EcdsaPrivateKey, domain string, resolverAddress crypto.EcdsaAddress) (string, error)
	HashDomain(input string) [32]byte
}

type IRegistryContract interface {
	Resolver(opts *bind.CallOpts, domain [32]byte) (crypto.EcdsaAddress, error)
	Register(opts *bind.TransactOpts, domain [32]byte, resolverAddress crypto.EcdsaAddress) (*types.Transaction, error)
}
