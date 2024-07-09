package commands

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

// register a provider
func RegisterProvider(cfg *config.Config, domain string, resolverAddress string, adminPrivateKey string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	defer ethClient.Close()
	registryService := services.NewRegistryService(cfg, ethClient)

	// recover the private key
	privateKey := crypto.EcdsaHexToPrivateKey(adminPrivateKey)

	// register the provider
	tx, err := registryService.RegisterDomain(privateKey, domain, common.HexToAddress(resolverAddress))
	if err != nil {
		return err
	}
	fmt.Printf("\nregistered provider %s. Transaction hash: %s\n", domain, tx)

	return nil
}
