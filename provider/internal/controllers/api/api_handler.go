package api_controller

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	api_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/api/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	crypto_service "github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	zip_service "github.com/matzapata/ipfs-compute/provider/pkg/zip"
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
	IPFS_PINATA_APIKEY := os.Getenv("IPFS_PINATA_APIKEY")
	IPFS_PINATA_SECRET := os.Getenv("IPFS_PINATA_SECRET")
	IPFS_PINATA_ENDPOINT := os.Getenv("IPFS_PINATA_ENDPOINT")

	// services
	cryptoRsaService := crypto_service.NewCryptoRsaService()
	cryptoEcdsaService := crypto_service.NewCryptoEcdsaService()
	zipService := zip_service.NewZipService()

	// constants
	providerEcdsaPrivateKey, err := cryptoEcdsaService.LoadPrivateKeyFromString(PROVIDER_ECDSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load ecdsa private key", err)
	}
	providerEcdsaAddress, err := cryptoEcdsaService.PrivateKeyToAddress(providerEcdsaPrivateKey)
	if err != nil {
		log.Fatal("cannot load ecdsa address ", err)
	}
	providerRsaPrivateKey, err := cryptoRsaService.LoadPrivateKeyFromString(PROVIDER_RSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load rsa private key", err)
	}
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal("cannot connect to rpc", err)
	}
	providerUnitPrice, _ := big.NewInt(0).SetString(PROVIDER_UNIT_PRICE, 10)

	// repositories
	artifactRepository := repositories.NewIpfsArtifactRepository(IPFS_PINATA_ENDPOINT, IPFS_PINATA_APIKEY, IPFS_PINATA_SECRET)

	// core services
	artifactsService := services.NewArtifactService(artifactRepository, cryptoRsaService, zipService)
	escrowService := services.NewEscrowService(ethClient, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)
	computeService := services.NewComputeService(
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
