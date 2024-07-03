package provider_controller

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/provider/routers"
	"github.com/matzapata/ipfs-compute/provider/internal/repositories"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type ApiHandler struct {
	Router *chi.Mux
}

func NewApiHandler() (*ApiHandler, error) {
	// load config
	cfg := config.ReadConfig("")

	// provider ecdsa keypair
	providerEcdsaPrivateKey, err := crypto.EcdsaLoadPrivateKeyFromString(cfg.ProviderEcdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	providerEcdsaAddress, err := crypto.EcdsaPrivateKeyToAddress(providerEcdsaPrivateKey)
	if err != nil {
		return nil, err
	}

	// provider rsa keypair
	providerRsaPrivateKey, err := crypto.RsaLoadPrivateKeyFromString(cfg.ProviderRsaPrivateKey)
	if err != nil {
		return nil, err
	}
	providerRsaPublicKey, err := crypto.RsaLoadPublicKeyFromString(cfg.ProviderRsaPublicKey)
	if err != nil {
		return nil, err
	}

	// provider unit price
	providerUnitPrice, success := big.NewInt(0).SetString(cfg.ProviderComputeUnitPrice, 10)
	if !success {
		return nil, err
	}

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
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	artifactsService := services.NewArtifactService(artifactRepository, unzipper)
	escrowService := services.NewEscrowService(ethClient, ethAuthenticator)
	computeService := services.NewComputeService(
		artifactsService,
		escrowService,
		providerEcdsaPrivateKey,
		&providerEcdsaAddress,
		providerRsaPrivateKey,
		providerRsaPublicKey,
		providerUnitPrice,
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
