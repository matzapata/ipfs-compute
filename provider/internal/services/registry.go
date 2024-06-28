package services

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	eth_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/eth"
	"golang.org/x/crypto/sha3"
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
	resolverAddress, err := r.Registry.Resolver(nil, r.HashDomain(domain))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Resolver address for domain %s: %s\n", fmt.Sprintf("%x", r.HashDomain(domain)), resolverAddress.Hex())

	if resolverAddress == common.HexToAddress("0x000000000000000000000000000000000000") {
		return nil, fmt.Errorf("resolver not found for domain %s", domain)
	}

	return contracts.NewResolver(resolverAddress, r.EthClient)
}

func (r *RegistryService) RegisterDomain(privateKey *ecdsa.PrivateKey, domain string, resolverAddress common.Address) (string, error) {
	auth, err := eth_helpers.BuildAuth(privateKey, r.EthClient, nil)
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
