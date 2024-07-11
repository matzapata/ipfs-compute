package commands

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

func BuildCommand(cfg *config.Config, serviceName string) error {
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return err
	}

	sourceRepo := repositories.NewSystemSourceRepository()
	artifactRepo := repositories.NewIpfsArtifactRepository(cfg.IpfsGateway, cfg.IpfsPinataApikey, cfg.IpfsPinataSecret)
	sourceService := services.NewSourceService(sourceRepo)
	registryService := services.NewRegistryService(cfg, ethClient)
	artifactBuilder := services.NewArtifactBuilderService(cfg, sourceService, registryService, artifactRepo)

	startTime := time.Now()
	err = artifactBuilder.BuildArtifact(serviceName)
	if err != nil {
		return err
	}

	fmt.Printf("\nBuilt in %v\n", time.Since(startTime))
	return nil
}
