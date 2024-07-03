package commands

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/console"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

// Deposit funds into the escrow account
func DepositCommand(cfg *config.Config, amount string, adminPrivateKey string) {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close()
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	escrowService := services.NewEscrowService(ethClient, ethAuthenticator)

	// confirm with the user
	prompt := fmt.Sprintf("You are about to deposit %s USDC into the escrow account. Continue? (y/n): ", amount)
	if !console.Confirm(prompt) {
		return
	}

	// recover private key
	privateKey, err := crypto.EcdsaLoadPrivateKeyFromString(adminPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// parse amount to bignumber
	bigAmount, success := new(big.Int).SetString(amount, 10)
	if !success {
		log.Fatal("couldn't parse amount")
	}

	// Approve escrow contract to spend USDC
	approve, err := escrowService.ApproveEscrow(privateKey, bigAmount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Approved escrow contract to spend %s USDC. Transaction hash: %s\n", amount, approve)

	// Deposit funds
	hash, err := escrowService.Deposit(privateKey, bigAmount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposit successful. Transaction hash: %s\n", hash)
}
