package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

// register a provider
func RegisterProvider(cfg *config.Config, domain string, resolverAddress string, adminPrivateKey string) {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close()
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	registryService := services.NewRegistryService(ethClient, ethAuthenticator)

	// recover the private key
	privateKey, err := crypto.EcdsaLoadPrivateKeyFromString(adminPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// register the provider
	tx, err := registryService.RegisterDomain(privateKey, domain, common.HexToAddress(resolverAddress))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered provider %s. Transaction hash: %s\n", domain, tx)

}
