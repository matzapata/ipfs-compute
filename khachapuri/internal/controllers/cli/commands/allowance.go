package commands

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

func AllowanceCommand(cfg *config.Config, adminAddr string, providerDomain string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	escrowService := services.NewEscrowService(cfg, ethClient)
	registryService := services.NewRegistryService(cfg, ethClient)

	// resolve domain
	providerDomainData, err := registryService.ResolveDomain(providerDomain)
	if err != nil {
		return err
	}

	// get allowance
	allowance, price, err := escrowService.Allowance(common.HexToAddress(adminAddr), common.HexToAddress(providerDomainData.EcdsaAddress))
	if err != nil {
		return err
	}
	fmt.Printf("\n%s allowed %s to consume %s USDC at %s per request\n", adminAddr, providerDomain, allowance.String(), price.String())

	return nil
}
