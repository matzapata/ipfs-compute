package services

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/matzapata/ipfs-compute/shared/contracts"
	registry "github.com/matzapata/ipfs-compute/shared/registry"
)

type RegistryService struct {
	EthClient *ethclient.Client
	Registry  *contracts.Registry
}

func NewRegistryService(client *ethclient.Client, registryAddress string) *RegistryService {
	address := common.HexToAddress(registryAddress)
	registry, err := contracts.NewRegistry(address, client)
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
	resolverAddress, err := r.Registry.Resolver(nil, registry.HashDomain(domain))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Resolver address for domain %s: %s\n", fmt.Sprintf("%x", registry.HashDomain(domain)), resolverAddress.Hex())

	if resolverAddress == common.HexToAddress("0x000000000000000000000000000000000000") {
		return nil, fmt.Errorf("resolver not found for domain %s", domain)
	}

	return contracts.NewResolver(resolverAddress, r.EthClient)
}

func (r *RegistryService) RegisterDomain(hexPrivateKey string, domain string, resolverAddress string) (string, error) {
	auth, err := buildAuth(hexPrivateKey, r.EthClient)
	if err != nil {
		return "", err
	}

	// Register domain
	tx, err := r.Registry.Register(auth, registry.HashDomain(domain), common.HexToAddress(resolverAddress))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
