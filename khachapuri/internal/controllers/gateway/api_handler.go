package gateway_controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	gateway_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type ApiHandler struct {
	EthClient ethclient.Client
}

func NewApiHandler() (*ApiHandler, error) {
	// load env vars
	RPC_URL := os.Getenv("RPC_URL")

	// create eth client
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return nil, err
	}

	return &ApiHandler{
		EthClient: *ethClient,
	}, nil
}

func (c *ApiHandler) Handle() {
	// create registry service
	ethAuthenticator := eth.NewEthAuthenticator(&c.EthClient)
	registryService := services.NewRegistryService(&c.EthClient, ethAuthenticator)

	// create router
	router := chi.NewRouter()

	// setup routes
	gateway_routers.SetupProxyRouter(router, registryService)

	// start server
	fmt.Println("Starting server on localhost:4000")
	err := http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}
}
