package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

func BalanceCommand(cfg *config.Config, address string) {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close()
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	escrowService := services.NewEscrowService(ethClient, ethAuthenticator)

	// get balance
	balance, err := escrowService.Balance(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance for address %s: %s\n", address, balance.String())
}
