package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

func ResolveCommand(cfg *config.Config, domain string) {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close()
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	registryService := services.NewRegistryService(ethClient, ethAuthenticator)

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
