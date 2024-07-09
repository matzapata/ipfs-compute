package commands

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

func ApproveCommand(cfg *config.Config, amount string, price string, providerDomain string, adminPrivateKey string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}
	escrowService := services.NewEscrowService(cfg, ethClient)
	registryService := services.NewRegistryService(cfg, ethClient)

	// resolve domain
	providerDomainData, err := registryService.ResolveDomain(providerDomain)
	if err != nil {
		return err
	}

	// parse user input
	adminPk := crypto.EcdsaHexToPrivateKey(adminPrivateKey)
	providerAddress := common.HexToAddress(providerDomainData.EcdsaAddress)
	amountInt := new(big.Int)
	amountInt, ok := amountInt.SetString(amount, 10)
	if !ok {
		return errors.New("invalid amount")
	}
	priceInt := new(big.Int)
	priceInt, ok = priceInt.SetString(price, 10)
	if !ok {
		return errors.New("invalid price")
	}

	// approve
	tx, err := escrowService.ApproveProvider(adminPk, providerAddress, amountInt, priceInt)
	if err != nil {
		return err
	}
	fmt.Printf("\napproved %s to consume %s USDC at %s per request\n", providerDomain, amount, price)
	fmt.Printf("tx: %s\n", tx)

	return nil
}
