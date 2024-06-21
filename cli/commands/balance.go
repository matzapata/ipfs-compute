package commands

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/services"
)

func BalanceCommand(address string, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}
	escrowService := services.NewEscrowService(ethclient, config.ESCROW_ADDRESS, config.USDC_ADDRESS)

	// get balance
	balance, err := escrowService.Balance(address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Balance for address %s: %s\n", address, balance.String())
}
