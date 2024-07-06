package commands

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/console"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

// Deposit funds into the escrow account
func DepositCommand(cfg *config.Config, amount string, adminPrivateKey string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	escrowService := services.NewEscrowService(ethClient, *cfg.EscrowAddress, *cfg.UsdcAddress)

	// confirm with the user
	prompt := fmt.Sprintf("You are about to deposit %s USDC into the escrow account. Continue? (y/n): ", amount)
	if !console.Confirm(prompt) {
		return errors.New("deposit cancelled")
	}

	// parse user input
	privateKey := crypto.EcdsaHexToPrivateKey(adminPrivateKey)
	bigAmount, success := new(big.Int).SetString(amount, 10)
	if !success {
		return errors.New("invalid amount")
	}

	// Deposit funds
	hashes, err := escrowService.Deposit(privateKey, bigAmount)
	if err != nil {
		return err
	}
	fmt.Printf("\nDeposit successful. \napproval hash: %s\ntransaction hash: %s\n", hashes.Approval, hashes.Deposit)

	return nil
}
