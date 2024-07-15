package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	provider_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/provider"
)

func main() {
	godotenv.Load(".env.provider")

	cachePath, err := os.MkdirTemp("cache", "khachapuri-cache-*")
	if err != nil {
		panic(err)
	}

	loader := config.NewEnvLoader()
	cfg := config.Config{
		RegistryAddress: config.RegistryAddress,
		EscrowAddress:   config.EscrowAddress,
		UsdcAddress:     config.UsdcAddress,
		ArtifactMaxSize: config.ArtifactMaxSize,
		TempPath:        os.TempDir(),
		CachePath:       cachePath,

		ProviderComputeUnitPrice: loader.LoadBigInt("PROVIDER_COMPUTE_UNIT_PRICE", true),
		EthRpc:                   loader.LoadString("ETH_RPC", true),
		IpfsGateway:              loader.LoadString("IPFS_GATEWAY", true),
		IpfsPinataApikey:         loader.LoadString("IPFS_PINATA_APIKEY", true),
		IpfsPinataSecret:         loader.LoadString("IPFS_PINATA_SECRET", true),
		ProviderEcdsaAddress:     loader.LoadEcdsaAddress("PROVIDER_ECDSA_ADDRESS", true),
		ProviderEcdsaPrivateKey:  loader.LoadEcdsaPrivateKey("PROVIDER_ECDSA_PRIVATE_KEY", true),
		ProviderRsaPrivateKey:    loader.LoadRsaPrivateKey("PROVIDER_RSA_PRIVATE_KEY", true),
		ProviderRsaPublicKey:     loader.LoadRsaPublicKey("PROVIDER_RSA_PUBLIC_KEY", true),
	}

	controller, err := provider_controller.NewApiHandler(&cfg)
	if err != nil {
		log.Panic(err)
	}

	controller.Handle(":4000")
}
