package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/services"
)

func ApproveCommand(privateKey string, rpc string, amount uint, price uint, providerDomain string) {
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

	// confirm with the user
	fmt.Printf("You are about to approve provider %s to spend %d USDC at %d per request. Continue? (y/n): ", providerDomain, amount, price)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		return
	}

	// Approve escrow contract to spend USDC
	approve, err := escrowService.ApproveProvider(privateKey, providerAddress.Hex(), amount, price)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Approved provider %s to spend %d USDC at %d per request. Transaction hash: %s\n", providerDomain, amount, price, approve)
}
