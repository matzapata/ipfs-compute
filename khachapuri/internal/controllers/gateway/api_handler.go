package gateway_controller

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

type ApiHandler struct {
	Router *chi.Mux
}

func NewApiHandler(cfg *config.Config) (*ApiHandler, error) {

	// create eth client
	ethClient, err := ethclient.Dial(cfg.EthRpc)
	if err != nil {
		return nil, err
	}

	// create registry service
	registryService := services.NewRegistryService(ethClient, *cfg.RegistryAddress)

	// create router
	router := chi.NewRouter()

	// setup routes
	routers.SetupProxyRouter(router, registryService)

	return &ApiHandler{
		Router: router,
	}, nil
}

func (c *ApiHandler) Handle(addr string) {
	fmt.Println("Starting server...")
	err := http.ListenAndServe(addr, c.Router)
	if err != nil {
		panic(err)
	}
}
