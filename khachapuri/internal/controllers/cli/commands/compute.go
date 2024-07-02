package commands

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
)

func ComputeCommand(config config.Config, cid string) {
	// eth connection
	eth, err := ethclient.Dial(config.Rpc)
	if err != nil {
		log.Fatal(err)
	}

	// repositories
	artifactRepository := repositories.NewIpfsArtifactRepository(config.IpfsGateway, config.IpfsApikey, config.IpfsSecret)

	// services

	zipService := zip_service.NewZipService()
	artifactsService := services.NewArtifactService(artifactRepository, zipService)
	escrowService := services.NewEscrowService(eth)
	computeService := services.NewComputeService(artifactsService, escrowService, config.ProviderEcdsaKey)

	// compute deployment

}
