package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/gateway/controllers"
	"github.com/matzapata/ipfs-compute/gateway/services"
	"github.com/matzapata/ipfs-compute/shared/registry"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Instantiate services
	proxyService := services.NewProxyService()
	registryService := registry.NewRegistryService(nil, common.HexToAddress("TODO:"))

	// Instantiate controllers
	proxyController := controllers.NewProxyController(registryService, proxyService)

	// register controllers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxyController.Proxy(w, r)
	})

	// Start the server
	log.Println("Starting proxy server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
