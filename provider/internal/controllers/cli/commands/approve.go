package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	console_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/console"
)

func ApproveCommand(hexPrivateKey string, rpc string, amount uint, price uint, providerDomain string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	registryService := services.NewRegistryService(ethclient, config.REGISTRY_ADDRESS)
	escrowService := services.NewEscrowService(ethclient, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)

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
	prompt := fmt.Sprintf("You are about to approve provider %s to spend %d USDC at %d per request. Continue? (y/n): ", providerDomain, amount, price)
	if !console_helpers.Confirm(prompt) {
		return
	}

	// recover private key
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Approve escrow contract to spend USDC
	approve, err := escrowService.ApproveProvider(privateKey, providerAddress.Hex(), amount, price)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Approved provider %s to spend %d USDC at %d per request. Transaction hash: %s\n", providerDomain, amount, price, approve)
}
