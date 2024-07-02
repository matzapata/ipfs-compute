package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

func BalanceCommand(config config.Config, address string) {
	eth, err := ethclient.Dial(config.Rpc)
	if err != nil {
		log.Fatal(err)
	}
	escrowService := services.NewEscrowService(eth)

	// get balance
	balance, err := escrowService.Balance(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance for address %s: %s\n", address, balance.String())
}
