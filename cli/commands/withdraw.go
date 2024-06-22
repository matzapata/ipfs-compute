package commands

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/helpers"
	"github.com/matzapata/ipfs-compute/shared/cryptoecdsa"
	"github.com/matzapata/ipfs-compute/shared/escrow"
)

// Withdraw funds from the escrow account
func WithdrawCommand(hexPrivateKey string, amount uint, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	escrowService := escrow.NewEscrowService(ethclient, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)

	// confirm with the user
	prompt := fmt.Sprintf("You are about to withdraw %d USDC from the escrow account. Continue? (y/n): ", amount)
	if !helpers.Confirm(prompt) {
		return
	}

	// recover private key
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Check balance
	address, err := cryptoecdsa.PrivateKeyToAddress(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := escrowService.Balance(address.Hex())
	if err != nil {
		log.Fatal(err)
	}
	if balance.Cmp(big.NewInt(int64(amount))) == -1 {
		log.Fatalf("Insufficient balance. Current balance: %s\n", balance.String())
	}

	// Withdraw funds
	hash, err := escrowService.Withdraw(privateKey, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Withdraw successful. Transaction hash: %s\n", hash)
}
