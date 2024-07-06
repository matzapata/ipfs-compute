package config

import (
	"math/big"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/spf13/viper"
)

type ViperLoader struct {
}

// TODO: receive defaults here
func NewViperLoader(configPath string) *ViperLoader {
	if configPath == "" {
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// read in environment variables that match. Usage:  `host: ${DB_HOST}`
	viper.AutomaticEnv()

	// set defaults
	viper.SetDefault("eth.usdc", "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	viper.SetDefault("eth.registry", "0xdb42A86B1bfe04E75B2A5F2bF7a3BBB52D7FFD2F")
	viper.SetDefault("eth.escrow", "0x5Fe8861F6571174a9564365384AE9b01CcdCd8D6")

	return &ViperLoader{}
}

func (e *ViperLoader) LoadString(key string, panicOnMissing bool) string {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return viper.GetString(key)
}

func (e *ViperLoader) LoadInt(key string, panicOnMissing bool) int {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return viper.GetInt(key)
}

func (e *ViperLoader) LoadUint(key string, panicOnMissing bool) uint {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return viper.GetUint(key)
}

func (e *ViperLoader) LoadBool(key string, panicOnMissing bool) bool {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return viper.GetBool(key)

}

func (e *ViperLoader) LoadBigInt(key string, panicOnMissing bool) *big.Int {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	bigInt := new(big.Int)
	bigInt, ok := bigInt.SetString(viper.GetString(key), 10)
	if !ok {
		panic("Error parsing string to big int")
	}

	return bigInt
}

func (e *ViperLoader) LoadEcdsaAddress(key string, panicOnMissing bool) *crypto.EcdsaAddress {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return crypto.EcdsaHexToAddress(viper.GetString(key))
}

func (e *ViperLoader) LoadEcdsaPrivateKey(key string, panicOnMissing bool) *crypto.EcdsaPrivateKey {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return crypto.EcdsaHexToPrivateKey(viper.GetString(key))
}

func (e *ViperLoader) LoadRsaPrivateKey(key string, panicOnMissing bool) *crypto.RsaPrivateKey {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return crypto.RsaLoadPrivateKeyFromString(viper.GetString(key))
}

func (e *ViperLoader) LoadRsaPublicKey(key string, panicOnMissing bool) *crypto.RsaPublicKey {
	if !viper.IsSet(key) && panicOnMissing {
		panic("missing var: " + key)
	}

	return crypto.RsaLoadPublicKeyFromString(viper.GetString(key))
}
