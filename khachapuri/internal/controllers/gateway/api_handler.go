package gateway_controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type ApiHandler struct {
	Router *chi.Mux
}

func NewApiHandler() (*ApiHandler, error) {
	// load env vars
	RPC_URL := os.Getenv("RPC_URL")

	// create eth client
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return nil, err
	}

	// create registry service
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	registryService := services.NewRegistryService(ethClient, ethAuthenticator)

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
