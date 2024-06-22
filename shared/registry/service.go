package registry

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/matzapata/ipfs-compute/shared/contracts"
	"github.com/matzapata/ipfs-compute/shared/eth"
)

type RegistryService struct {
	EthClient *ethclient.Client
	Registry  *contracts.Registry
}

func NewRegistryService(client *ethclient.Client, registryAddress common.Address) *RegistryService {
	registry, err := contracts.NewRegistry(registryAddress, client)
	if err != nil {
		panic(err)
	}

	return &RegistryService{
		EthClient: client,
		Registry:  registry,
	}
}

// Resolve the domain to get the provider
func (r *RegistryService) ResolveDomain(domain string) (*contracts.Resolver, error) {
	// get resolver address and instantiate it
	resolverAddress, err := r.Registry.Resolver(nil, HashDomain(domain))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Resolver address for domain %s: %s\n", fmt.Sprintf("%x", HashDomain(domain)), resolverAddress.Hex())

	if resolverAddress == common.HexToAddress("0x000000000000000000000000000000000000") {
		return nil, fmt.Errorf("resolver not found for domain %s", domain)
	}

	return contracts.NewResolver(resolverAddress, r.EthClient)
}

func (r *RegistryService) RegisterDomain(privateKey *ecdsa.PrivateKey, domain string, resolverAddress common.Address) (string, error) {
	auth, err := eth.BuildAuth(privateKey, r.EthClient, nil)
	if err != nil {
		return "", err
	}

	// Register domain
	tx, err := r.Registry.Register(auth, HashDomain(domain), resolverAddress)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
