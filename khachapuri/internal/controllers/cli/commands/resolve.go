package commands

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

func ResolveCommand(cfg *config.Config, domain string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	registryService := services.NewRegistryService(cfg, ethClient)

	// resolve domain
	providerDomainData, err := registryService.ResolveDomain(domain)
	if err != nil {
		return err
	}

	fmt.Println("\nProvider address:", providerDomainData.EcdsaAddress)
	fmt.Println("Provider public key:", providerDomainData.RsaPublicKey)
	fmt.Println("Provider server endpoint:", providerDomainData.ServerEndpoint)

	return nil
}
