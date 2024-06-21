package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/services"
)

// Deposit funds into the escrow account
func DepositCommand(privateKey string, amount uint, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	escrowService := services.NewEscrowService(ethclient, config.ESCROW_ADDRESS, config.USDC_ADDRESS)

	// confirm with the user
	fmt.Printf("You are about to deposit %d USDC into the escrow account. Continue? (y/n): ", amount)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		return
	}

	// Approve escrow contract to spend USDC
	approve, err := escrowService.ApproveEscrow(privateKey, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Approved escrow contract to spend %d USDC. Transaction hash: %s\n", amount, approve)

	// Deposit funds
	hash, err := escrowService.Deposit(privateKey, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposit successful. Transaction hash: %s\n", hash)
}
