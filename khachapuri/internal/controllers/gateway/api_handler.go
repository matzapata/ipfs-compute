package gateway_controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	gateway_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (c *ApiHandler) Handle() {
	// load env vars
	RPC_URL := os.Getenv("RPC_URL")

	// create eth client
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal("cannot connect to rpc", err)
	}

	// create registry service
	registryService := services.NewRegistryService(ethClient, config.REGISTRY_ADDRESS)

	// create router
	router := chi.NewRouter()

	// setup routes
	gateway_routers.SetupProxyRouter(router, registryService)

	// start server
	fmt.Println("Starting server on localhost:4000")
	err = http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}
}
