package provider_controller

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/provider/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
)

type ApiHandler struct {
	Router *chi.Mux
}

func NewApiHandler(cfg *config.Config) (*ApiHandler, error) {
	// eth client
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return nil, err
	}

	// repositories
	artifactRepository := repositories.NewIpfsArtifactRepository(
		cfg.IpfsGateway,
		cfg.IpfsPinataApikey,
		cfg.IpfsPinataSecret,
	)

	// core services
	unzipper := archive.NewUnzipper()
	artifactsService := services.NewArtifactService(artifactRepository, unzipper, cfg.ArtifactMaxSize)
	escrowService := services.NewEscrowService(ethClient, *cfg.EscrowAddress, *cfg.UsdcAddress)
	computeService := services.NewComputeService(
		artifactsService,
		escrowService,
		cfg.ProviderEcdsaPrivateKey,
		cfg.ProviderEcdsaAddress,
		cfg.ProviderRsaPrivateKey,
		cfg.ProviderRsaPublicKey,
		cfg.ProviderComputeUnitPrice,
	)

	// router
	router := chi.NewRouter()

	// setup routes
	routers.SetupComputeRoutes(router, computeService)
	routers.SetupHealthRotes(router)

	return &ApiHandler{
		Router: router,
	}, nil
}

func (a *ApiHandler) Handle(addr string) {
	fmt.Println("Starting server...")
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		panic(err)
	}
}
