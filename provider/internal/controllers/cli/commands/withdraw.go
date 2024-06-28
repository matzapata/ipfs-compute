package commands

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	console_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/console"
)

// Withdraw funds from the escrow account
func WithdrawCommand(hexPrivateKey string, amount uint, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	escrowService := services.NewEscrowService(ethclient, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)
	cryptoEcdsaService := crypto_service.NewCryptoEcdsaService()

	// confirm with the user
	prompt := fmt.Sprintf("You are about to withdraw %d USDC from the escrow account. Continue? (y/n): ", amount)
	if !console_helpers.Confirm(prompt) {
		return
	}

	// recover private key
	privateKey, err := cryptoEcdsaService.LoadPrivateKeyFromString(hexPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Check balance
	address, err := cryptoEcdsaService.PrivateKeyToAddress(privateKey)
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
