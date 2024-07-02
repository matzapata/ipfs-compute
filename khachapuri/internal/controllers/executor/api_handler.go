package executor_controller

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	executor_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/executor/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type ApiHandler struct {
	ProviderEcdsaAddress     crypto.EcdsaAddress
	ProviderEcdsaPrivateKey  crypto.EcdsaPrivateKey
	ProviderRsaPrivateKey    crypto.RsaPrivateKey
	ProviderRsaPublicKey     crypto.RsaPublicKey
	ProviderComputeUnitPrice *big.Int
	EthClient                ethclient.Client
	IpfsGateway              string
	IpfsPinataApikey         string
	IpfsPinataSecret         string
}

func NewApiHandler() (*ApiHandler, error) {
	// load environment variables and configs
	RPC_URL := os.Getenv("RPC_URL")
	PROVIDER_ECDSA_PRIVATE_KEY := os.Getenv("PROVIDER_ECDSA_PRIVATE_KEY")
	PROVIDER_RSA_PRIVATE_KEY := os.Getenv("PROVIDER_RSA_PRIVATE_KEY")
	PROVIDER_RSA_PUBLIC_KEY := os.Getenv("PROVIDER_RSA_PUBLIC_KEY")
	PROVIDER_UNIT_PRICE := os.Getenv("PROVIDER_UNIT_PRICE")
	IPFS_PINATA_APIKEY := os.Getenv("IPFS_PINATA_APIKEY")
	IPFS_PINATA_SECRET := os.Getenv("IPFS_PINATA_SECRET")
	IPFS_PINATA_ENDPOINT := os.Getenv("IPFS_PINATA_ENDPOINT")

	// constants
	providerEcdsaPrivateKey, err := crypto.EcdsaLoadPrivateKeyFromString(PROVIDER_ECDSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load ecdsa private key", err)
	}
	providerEcdsaAddress, err := crypto.EcdsaPrivateKeyToAddress(providerEcdsaPrivateKey)
	if err != nil {
		log.Fatal("cannot load ecdsa address ", err)
	}
	providerRsaPrivateKey, err := crypto.RsaLoadPrivateKeyFromString(PROVIDER_RSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal("cannot load rsa private key", err)
	}
	providerRsaPublicKey, err := crypto.RsaLoadPublicKeyFromString(PROVIDER_RSA_PUBLIC_KEY)
	if err != nil {
		log.Fatal("cannot load rsa public key", err)
	}
	ethClient, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal("cannot connect to rpc", err)
	}
	providerUnitPrice, success := big.NewInt(0).SetString(PROVIDER_UNIT_PRICE, 10)
	if !success {
		log.Fatal("couldn't parse provider unit price")
	}

	return &ApiHandler{
		ProviderEcdsaAddress:     providerEcdsaAddress,
		ProviderEcdsaPrivateKey:  *providerEcdsaPrivateKey,
		ProviderRsaPrivateKey:    *providerRsaPrivateKey,
		ProviderRsaPublicKey:     *providerRsaPublicKey,
		ProviderComputeUnitPrice: providerUnitPrice,
		EthClient:                *ethClient,
		IpfsGateway:              IPFS_PINATA_ENDPOINT,
		IpfsPinataApikey:         IPFS_PINATA_APIKEY,
		IpfsPinataSecret:         IPFS_PINATA_SECRET,
	}, nil
}

func (a *ApiHandler) Handle() {
	// repositories
	artifactRepository := repositories.NewIpfsArtifactRepository(a.IpfsGateway, a.IpfsPinataApikey, a.IpfsPinataSecret)

	// core services
	unzipper := archive.NewUnzipper()
	ethAuthenticator := eth.NewEthAuthenticator(&a.EthClient)
	artifactsService := services.NewArtifactService(artifactRepository, unzipper)
	escrowService := services.NewEscrowService(&a.EthClient, ethAuthenticator)
	computeService := services.NewComputeService(
		artifactsService,
		escrowService,
		&a.ProviderEcdsaPrivateKey,
		&a.ProviderEcdsaAddress,
		&a.ProviderRsaPrivateKey,
		&a.ProviderRsaPublicKey,
		a.ProviderComputeUnitPrice,
	)

	// create router
	router := chi.NewRouter()

	// setup routes
	executor_routers.SetupComputeRoutes(router, computeService)
	executor_routers.SetupHealthRotes(router)

	// start server
	fmt.Println("Starting server on localhost:4000")
	err := http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}
}
