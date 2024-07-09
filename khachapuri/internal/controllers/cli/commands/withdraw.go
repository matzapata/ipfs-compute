package commands

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

// Withdraw funds from the escrow account
func WithdrawCommand(cfg *config.Config, amount string, adminPrivateKey string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	escrowService := services.NewEscrowService(cfg, ethClient)

	// parse user input
	privateKey := crypto.EcdsaHexToPrivateKey(adminPrivateKey)
	amountInt := new(big.Int)
	amountInt, ok := amountInt.SetString(amount, 10)
	if !ok {
		return errors.New("invalid amount")
	}

	// Withdraw funds
	tx, err := escrowService.Withdraw(privateKey, amountInt)
	if err != nil {
		return err
	}

	fmt.Printf("\nWithdrawal successful.\nTransaction hash: %s\n", tx)

	return nil
}
