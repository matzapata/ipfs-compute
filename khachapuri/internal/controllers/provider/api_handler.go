package provider_controller

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/provider/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
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
	fileSystemCache := repositories.NewFileSystemCache(".cache")
	artifactRepository := repositories.NewIpfsArtifactRepository(
		fileSystemCache,
		cfg.IpfsGateway,
		cfg.IpfsPinataApikey,
		cfg.IpfsPinataSecret,
	)

	// core services
	artifactsService := services.NewArtifactService(cfg, artifactRepository)
	escrowService := services.NewEscrowService(cfg, ethClient)
	computeExecutor := services.NewComputeExecutor(cfg)
	computeService := services.NewComputeService(
		cfg,
		artifactsService,
		escrowService,
		computeExecutor,
	)

	// router
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// setup routes
	routers.SetupComputeRoutes(router, computeService)
	routers.SetupHealthRotes(router)

	return &ApiHandler{
		Router: router,
	}, nil
}

func (a *ApiHandler) Handle(addr string) {
	fmt.Println("Starting server at localhost" + addr)
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		panic(err)
	}
}
