package main

import (
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/controllers"
	"github.com/matzapata/ipfs-compute/repositories"
	ipfsRepositories "github.com/matzapata/ipfs-compute/repositories/ipfs"
	services "github.com/matzapata/ipfs-compute/services"
)

func main() {
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
}
