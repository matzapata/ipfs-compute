package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/services"
)

// register a provider
func RegisterProvider(hexPrivateKey string, rpc string, domain string, resolverAddress string) {
	// dial the ethereum client
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	// create a new registry service
	registryService := services.NewRegistryService(ethclient, config.REGISTRY_ADDRESS)

	// register the provider
	tx, err := registryService.RegisterDomain(hexPrivateKey, domain, resolverAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered provider %s. Transaction hash: %s\n", domain, tx)

}
