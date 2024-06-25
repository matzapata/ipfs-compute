package api_controller

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/artifact"
	artifact_repository "github.com/matzapata/ipfs-compute/provider/internal/artifact/repository"
	"github.com/matzapata/ipfs-compute/provider/internal/compute"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	api_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/api/routers"
	"github.com/matzapata/ipfs-compute/provider/pkg/escrow"
	ecdsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/ecdsa"
	rsa_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/rsa"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (a *ApiHandler) Handle() {
	// load environment variables and configs
	RPC_URL := os.Getenv("RPC_URL")
	PROVIDER_ECDSA_PRIVATE_KEY := os.Getenv("PROVIDER_ECDSA_PRIVATE_KEY")
	PROVIDER_RSA_PRIVATE_KEY := os.Getenv("PROVIDER_RSA_PRIVATE_KEY")
	PROVIDER_UNIT_PRICE := os.Getenv("PROVIDER_UNIT_PRICE")

	// constants
	providerEcdsaPrivateKey, err := ecdsa_helpers.HexToPrivateKey(PROVIDER_ECDSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load ecdsa private key", err)
	}
	providerEcdsaAddress, err := ecdsa_helpers.PrivateKeyToAddress(providerEcdsaPrivateKey)
	if err != nil {
		log.Fatal("cannot load ecdsa address ", err)
	}
	providerRsaPrivateKey, err := rsa_helpers.LoadPrivateKeyFromString(PROVIDER_RSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load rsa private key", err)
	}
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal("cannot connect to rpc", err)
	}
	providerUnitPrice, _ := big.NewInt(0).SetString(PROVIDER_UNIT_PRICE, 10)

	// instantiate repositories
	artifactRepository := artifact_repository.NewIpfsArtifactRepository()

	// instantiate services
	artifactsService := artifact.NewArtifactService(artifactRepository)
	escrowService := escrow.NewEscrowService(ethClient, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)
	computeService := compute.NewComputeService(
		artifactsService,
		escrowService,
		providerEcdsaPrivateKey,
		&providerEcdsaAddress,
		providerRsaPrivateKey,
		providerUnitPrice,
	)

	// create router
	router := chi.NewRouter()

	// setup routes
	api_routers.SetupComputeRoutes(router, computeService)
	api_routers.SetupHealthRotes(router)

	// start server
	fmt.Println("Starting server on localhost:4000")
	err = http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}
}
