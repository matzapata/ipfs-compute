package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

func AllowanceCommand(cfg *config.Config, adminAddr string, providerAddr string) {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close()
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	escrowService := services.NewEscrowService(ethClient, ethAuthenticator)

	// get allowance
	allowance, price, err := escrowService.Allowance(common.HexToAddress(adminAddr), common.HexToAddress(providerAddr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s allowed %s to consume %s USDC at %s per request\n", adminAddr, providerAddr, allowance.String(), price.String())
}
