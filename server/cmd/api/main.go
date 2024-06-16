package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/controllers"
	"github.com/matzapata/ipfs-compute/repositories"
	ipfsRepositories "github.com/matzapata/ipfs-compute/repositories/ipfs"
	services "github.com/matzapata/ipfs-compute/services"
)

func main() {
	// load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// repositories
	var deploymentRepository repositories.DeploymentsRepository = ipfsRepositories.NewIpfsDeploymentsRepository()

	// services
	computeService := services.NewComputeService()
	deploymentService := services.NewDeploymentService(&deploymentRepository)

	// controllers
	computeController := controllers.NewComputeHttpController(computeService, deploymentService)

	// router
	router := chi.NewRouter()
	router.HandleFunc("/{cid}", computeController.Compute)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to IPFS Compute"))
	})

	// start server
	http.ListenAndServe(":3000", router)
}
