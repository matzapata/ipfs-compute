package services

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
	"golang.org/x/crypto/sha3"
)

type RegistryService struct {
	EthClient        *ethclient.Client
	Registry         domain.IRegistryContract
	EthAuthenticator eth.IEthAuthenticator
}

func NewRegistryService(ethClient *ethclient.Client, registryAddress crypto.EcdsaAddress) *RegistryService {
	registry, err := contracts.NewRegistry(registryAddress, ethClient)
	if err != nil {
		panic(err)
	}
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)

	return &RegistryService{
		EthClient:        ethClient,
		Registry:         registry,
		EthAuthenticator: ethAuthenticator,
	}
}

// Resolve the domain to get the provider
func (r *RegistryService) ResolveDomain(domainName string) (*domain.ProviderDomainData, error) {
	// get resolver address and instantiate it
	resolverAddress, err := r.Registry.Resolver(nil, r.HashDomain(domainName))
	if err != nil {
		return nil, err
	}

	if (resolverAddress == common.Address{}) {
		return nil, fmt.Errorf("resolver not found for domain %s", domainName)
	}

	resolver, err := contracts.NewResolver(resolverAddress, r.EthClient)
	if err != nil {
		return nil, err
	}

	// resolve all fields
	ecdsaAddress, err := resolver.Addr(nil)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, err := resolver.Pubkey(nil)
	if err != nil {
		return nil, err
	}
	server, err := resolver.Server(nil)
	if err != nil {
		return nil, err
	}

	return &domain.ProviderDomainData{
		EcdsaAddress:   ecdsaAddress.Hex(),
		RsaPublicKey:   rsaPublicKey,
		ServerEndpoint: server,
	}, nil
}

// Register the domain and returns the transaction hash
func (r *RegistryService) RegisterDomain(privateKey *crypto.EcdsaPrivateKey, domain string, resolverAddress crypto.EcdsaAddress) (string, error) {
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

// Ensures consistent hashing of domain for storage
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
