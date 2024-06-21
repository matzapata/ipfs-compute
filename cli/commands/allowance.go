package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/services"
)

func AllowanceCommand(address string, providerDomain string, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	registryService := services.NewRegistryService(ethclient, config.REGISTRY_ADDRESS)
	escrowService := services.NewEscrowService(ethclient, config.ESCROW_ADDRESS, config.USDC_ADDRESS)

	// resolve domain
	resolver, err := registryService.ResolveDomain(providerDomain)
	if err != nil {
		log.Fatal(err)
	}
	providerAddress, err := resolver.Addr(nil)
	if err != nil {
		log.Fatal(err)
	}

	// get allowance
	allowance, price, err := escrowService.Allowance(address, providerAddress.Hex())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s allowed %s to consume %s USDC at %s per request\n", address, providerDomain, allowance.String(), price.String())
}
