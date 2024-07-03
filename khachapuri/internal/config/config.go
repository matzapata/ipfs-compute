package config

import "github.com/spf13/viper"

type Config struct {
	EthRpc                   string
	ProviderEcdsaAddress     string
	ProviderEcdsaPrivateKey  string
	ProviderRsaPrivateKey    string
	ProviderRsaPublicKey     string
	ProviderComputeUnitPrice string
	IpfsGateway              string
	IpfsPinataApikey         string
	IpfsPinataSecret         string
}

func ReadConfig(configPath string) *Config {
	if configPath == "" {
		configPath = "config.yaml"
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil
	}

	return &Config{
		EthRpc:                   viper.GetString("eth.rpc"),
		ProviderEcdsaAddress:     viper.GetString("provider.ecdsa.address"),
		ProviderEcdsaPrivateKey:  viper.GetString("provider.ecdsa.private_key"),
		ProviderRsaPrivateKey:    viper.GetString("provider.rsa.private_key"),
		ProviderRsaPublicKey:     viper.GetString("provider.rsa.public_key"),
		ProviderComputeUnitPrice: viper.GetString("provider.compute.unit_price"),
		IpfsGateway:              viper.GetString("ipfs.gateway"),
		IpfsPinataApikey:         viper.GetString("ipfs.apikey"),
		IpfsPinataSecret:         viper.GetString("ipfs.secret"),
	}
}
