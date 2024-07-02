package services

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
	"golang.org/x/crypto/sha3"
)

type RegistryService struct {
	EthClient        *ethclient.Client
	Registry         *contracts.Registry
	EthAuthenticator eth.IEthAuthenticator
}

func NewRegistryService(client *ethclient.Client, ethAuthenticator eth.IEthAuthenticator) *RegistryService {
	registry, err := contracts.NewRegistry(config.REGISTRY_ADDRESS, client)
	if err != nil {
		panic(err)
	}

	return &RegistryService{
		EthClient:        client,
		Registry:         registry,
		EthAuthenticator: ethAuthenticator,
	}
}

// Resolve the domain to get the provider
func (r *RegistryService) ResolveDomain(domain string) (*contracts.Resolver, error) {
	// get resolver address and instantiate it
	resolverAddress, err := r.Registry.Resolver(nil, r.HashDomain(domain))
	if err != nil {
		return nil, err
	}

	if (resolverAddress == common.Address{}) {
		return nil, fmt.Errorf("resolver not found for domain %s", domain)
	}

	return contracts.NewResolver(resolverAddress, r.EthClient)
}

func (r *RegistryService) ResolveServer(domain string) (string, error) {
	resolver, err := r.ResolveDomain(domain)
	if err != nil {
		return "", err
	}

	return resolver.Server(nil)
}

func (r *RegistryService) RegisterDomain(privateKey *ecdsa.PrivateKey, domain string, resolverAddress common.Address) (string, error) {
	auth, err := r.EthAuthenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	// Register domain
	tx, err := r.Registry.Register(auth, r.HashDomain(domain), resolverAddress)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (r *RegistryService) HashDomain(input string) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	_, _ = hash.Write([]byte(input))

	// Get the resulting encoded byte slice
	sha3 := hash.Sum(nil)

	var sha32Bytes [32]byte
	copy(sha32Bytes[:], sha3)

	// Convert the encoded byte slice to a string
	return sha32Bytes
}
