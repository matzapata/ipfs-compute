package config

import (
	"math/big"
	"os"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type Config struct {
	IpfsGateway       string
	IpfsApikey        string
	IpfsSecret        string
	Rpc               string
	ProviderPriceUnit *big.Int
	ProviderRsaKey    *crypto.RsaPrivateKey
	ProviderEcdsaKey  *crypto.EcdsaPrivateKey
}

func LoadConfigFromEnv() *Config {
	// Load config from environment variables

	// Load keys
	providerRsaKey, _ := crypto.RsaLoadPrivateKeyFromString(os.Getenv("PROVIDER_RSA_KEY"))
	providerEcdsaKey, _ := crypto.EcdsaLoadPrivateKeyFromString(os.Getenv("PROVIDER_ECDSA_KEY"))

	// Load price unit
	providerPriceUnit := big.NewInt(0)
	providerPriceUnit.SetString(os.Getenv("PROVIDER_PRICE_UNIT"), 10)

	return &Config{
		IpfsGateway:       os.Getenv("IPFS_GATEWAY"),
		IpfsApikey:        os.Getenv("IPFS_APIKEY"),
		IpfsSecret:        os.Getenv("IPFS_SECRET"),
		Rpc:               os.Getenv("RPC"),
		ProviderPriceUnit: providerPriceUnit,
		ProviderRsaKey:    providerRsaKey,
		ProviderEcdsaKey:  providerEcdsaKey,
	}
}
