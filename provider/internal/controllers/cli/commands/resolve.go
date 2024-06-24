package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/pkg/registry"
)

func ResolveCommand(rpc string, domain string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	registryService := registry.NewRegistryService(ethclient, config.REGISTRY_ADDRESS)

	// resolve domain
	resolver, err := registryService.ResolveDomain(domain)
	if err != nil {
		log.Fatal(err)
	}

	// get provider
	address, err := resolver.Addr(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Provider address:", address.Hex())

	// get public key
	publicKey, err := resolver.Pubkey(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Provider public key:", publicKey)

	// get endpoint
	endpoint, err := resolver.Server(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Provider server endpoint:", endpoint)
}
