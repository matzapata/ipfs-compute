package main

import (
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/config"
	"github.com/matzapata/ipfs-compute/provider/controllers"
	"github.com/matzapata/ipfs-compute/provider/repositories"
	ipfsRepositories "github.com/matzapata/ipfs-compute/provider/repositories/ipfs"
	services "github.com/matzapata/ipfs-compute/provider/services"
	"github.com/matzapata/ipfs-compute/shared/escrow"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	RPC_URL := os.Getenv("RPC_URL")
	UNIT_PRICE := new(big.Int)
	UNIT_PRICE.SetString(os.Getenv("UNIT_PRICE"), 10)
	PROVIDER_ECDSA_PRIVATE_KEY := os.Getenv("PROVIDER_ECDSA_PRIVATE_KEY")

	// load private key
	providerEcdsaPrivateKey, err := crypto.HexToECDSA(PROVIDER_ECDSA_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}

	// eth client
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal(err)
	}

	// repositories
	var deploymentRepository repositories.DeploymentsRepository = ipfsRepositories.NewIpfsDeploymentsRepository()

	// services
	computeService := services.NewComputeService()
	deploymentService := services.NewDeploymentService(&deploymentRepository)
	escrowService := escrow.NewEscrowService(client, &config.ESCROW_ADDRESS, &config.USDC_ADDRESS)

	// controllers
	computeController := controllers.NewComputeHttpController(
		computeService,
		deploymentService,
		escrowService,
		UNIT_PRICE,
		providerEcdsaPrivateKey,
	)

	// router
	router := chi.NewRouter()
	router.HandleFunc("/{cid}", computeController.Compute)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// start server
	http.ListenAndServe(":8080", router)
}
