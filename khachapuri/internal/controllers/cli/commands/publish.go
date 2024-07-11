package commands

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/console"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

func PublishCommand(cfg *config.Config, serviceName string, providerDomain string, adminKey string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}

	sourceRepo := repositories.NewSystemSourceRepository()
	artifactRepo := repositories.NewIpfsArtifactRepository(cfg.IpfsGateway, cfg.IpfsPinataApikey, cfg.IpfsPinataSecret)
	sourceService := services.NewSourceService(sourceRepo)
	registryService := services.NewRegistryService(cfg, ethClient)
	artifactBuilder := services.NewArtifactBuilderService(cfg, sourceService, registryService, artifactRepo)

	// request confirmation
	prompt := fmt.Sprintf("You're making your code available and giving %v permanent access to your env variables. Continue? (y/n): ", providerDomain)
	if !console.Confirm(prompt) {
		return errors.New("deployment cancelled")
	}

	cid, err := artifactBuilder.PublishArtifact(serviceName, providerDomain, crypto.EcdsaHexToPrivateKey(adminKey))
	if err != nil {
		return err
	}
	fmt.Printf("Artifact published with CID: %v\n", cid)

	return nil
}
