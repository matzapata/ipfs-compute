package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	cli_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/cli"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// viperLoader := config.NewViperLoader(".")
	envLoader := config.NewEnvLoader()
	cfg := config.Config{
		RegistryAddress: config.RegistryAddress,
		EscrowAddress:   config.EscrowAddress,
		UsdcAddress:     config.UsdcAddress,
		ArtifactMaxSize: config.ArtifactMaxSize,

		EthRpc:                   envLoader.LoadString("ETH_RPC", true),
		IpfsGateway:              envLoader.LoadString("IPFS_GATEWAY", true),
		IpfsPinataApikey:         envLoader.LoadString("IPFS_PINATA_APIKEY", true),
		IpfsPinataSecret:         envLoader.LoadString("IPFS_PINATA_SECRET", true),
		ProviderEcdsaAddress:     envLoader.LoadEcdsaAddress("PROVIDER_ECDSA_ADDRESS", true),
		ProviderEcdsaPrivateKey:  envLoader.LoadEcdsaPrivateKey("PROVIDER_ECDSA_PRIVATE_KEY", true),
		ProviderRsaPrivateKey:    envLoader.LoadRsaPrivateKey("PROVIDER_RSA_PRIVATE_KEY", true),
		ProviderRsaPublicKey:     envLoader.LoadRsaPublicKey("PROVIDER_RSA_PUBLIC_KEY", true),
		ProviderComputeUnitPrice: envLoader.LoadBigInt("PROVIDER_COMPUTE_UNIT_PRICE", true),
	}

	controller := cli_controller.NewCliHandler(&cfg)
	controller.Handle()
}
