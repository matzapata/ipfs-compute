package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/server/controllers"
	"github.com/matzapata/ipfs-compute/server/repositories"
	ipfsRepositories "github.com/matzapata/ipfs-compute/server/repositories/ipfs"
	services "github.com/matzapata/ipfs-compute/server/services"
)

func main() {
	// load environment variables
	err := godotenv.Load()
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
		w.Write([]byte("OK"))
	})

	// start server
	http.ListenAndServe(":3000", router)
}
