package main

import (
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	gateway_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway"
)

func main() {
	godotenv.Load(".env.gateway")

	loader := config.NewEnvLoader()
	cfg := config.Config{
		RegistryAddress: config.RegistryAddress,
		EscrowAddress:   config.EscrowAddress,
		UsdcAddress:     config.UsdcAddress,
		ArtifactMaxSize: config.ArtifactMaxSize,

		EthRpc: loader.LoadString("ETH_RPC", true),
	}

	controller, err := gateway_controller.NewApiHandler(&cfg)
	if err != nil {
		panic(err)
	}

	controller.Handle(":3000")
}
