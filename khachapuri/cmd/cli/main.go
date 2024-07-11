package main

import (
	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	cli_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/cli"
)

func main() {
	godotenv.Load() // ignore errors, it's fine if there's no .env file. We check for missing env vars later.

	envLoader := config.NewEnvLoader()
	cfg := config.Config{
		RegistryAddress: config.RegistryAddress,
		EscrowAddress:   config.EscrowAddress,
		UsdcAddress:     config.UsdcAddress,
		ArtifactMaxSize: config.ArtifactMaxSize,
		BuildDir:        ".khachapuri",

		EthRpc:                   envLoader.LoadString("ETH_RPC", true),
		IpfsGateway:              envLoader.LoadString("IPFS_GATEWAY", true),
		IpfsPinataApikey:         envLoader.LoadString("IPFS_PINATA_APIKEY", true),
		IpfsPinataSecret:         envLoader.LoadString("IPFS_PINATA_SECRET", true),
		ProviderEcdsaAddress:     envLoader.LoadEcdsaAddress("PROVIDER_ECDSA_ADDRESS", false),
		ProviderEcdsaPrivateKey:  envLoader.LoadEcdsaPrivateKey("PROVIDER_ECDSA_PRIVATE_KEY", false),
		ProviderRsaPrivateKey:    envLoader.LoadRsaPrivateKey("PROVIDER_RSA_PRIVATE_KEY", false),
		ProviderRsaPublicKey:     envLoader.LoadRsaPublicKey("PROVIDER_RSA_PUBLIC_KEY", false),
		ProviderComputeUnitPrice: envLoader.LoadBigInt("PROVIDER_COMPUTE_UNIT_PRICE", false),
	}

	controller := cli_controller.NewCliHandler(&cfg)
	controller.Handle()
}
