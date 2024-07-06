package commands

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

func BalanceCommand(cfg *config.Config, address string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	escrowService := services.NewEscrowService(ethClient, *cfg.EscrowAddress, *cfg.UsdcAddress)

	// get balance
	balance, err := escrowService.Balance(common.HexToAddress(address))
	if err != nil {
		return err
	}
	fmt.Printf("\nbalance for address %s: %s\n", address, balance.String())

	return nil
}
