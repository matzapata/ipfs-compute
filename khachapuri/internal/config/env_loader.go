package config

import (
	"math/big"
	"os"
	"strconv"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type EnvLoader struct {
}

func NewEnvLoader() *EnvLoader {
	return &EnvLoader{}
}

func (e *EnvLoader) LoadString(key string, panicOnMissing bool) string {
	val := os.Getenv(key)
	if val == "" && panicOnMissing {
		panic("missing env var: " + key)
	}

	return val
}

func (e *EnvLoader) LoadInt(key string, panicOnMissing bool) int {
	val := os.Getenv(key)
	if val == "" && panicOnMissing {
		panic("missing env var: " + key)
	}

	number, err := strconv.Atoi(val)
	if err != nil {
		panic("Error parsing string to integer")

	}

	return number
}

func (e *EnvLoader) LoadUint(key string, panicOnMissing bool) uint {
	val := os.Getenv(key)
	if val == "" && panicOnMissing {
		panic("missing env var: " + key)
	}

	number, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic("Error parsing string to unsigned integer")
	}

	return uint(number)
}

func (e *EnvLoader) LoadBool(key string, panicOnMissing bool) bool {
	val := os.Getenv(key)
	if val == "" && panicOnMissing {
		panic("missing env var: " + key)
	}

	boolean, err := strconv.ParseBool(val)
	if err != nil {
		panic("Error parsing string to boolean")
	}

	return boolean

}

func (e *EnvLoader) LoadBigInt(key string, panicOnMissing bool) *big.Int {
	val := os.Getenv(key)
	if val == "" {
		if panicOnMissing {
			panic("missing env var: " + key)
		} else {
			return nil
		}
	}

	bigInt := new(big.Int)
	bigInt, ok := bigInt.SetString(val, 10)
	if !ok {
		panic("Error parsing string to big int")
	}

	return bigInt
}

func (e *EnvLoader) LoadEcdsaAddress(key string, panicOnMissing bool) *crypto.EcdsaAddress {
	val := os.Getenv(key)
	if val == "" {
		if panicOnMissing {
			panic("missing env var: " + key)
		} else {
			return nil
		}
	}

	return crypto.EcdsaHexToAddress(val)
}

func (e *EnvLoader) LoadEcdsaPrivateKey(key string, panicOnMissing bool) *crypto.EcdsaPrivateKey {
	val := os.Getenv(key)
	if val == "" {
		if panicOnMissing {
			panic("missing env var: " + key)
		} else {
			return nil
		}
	}

	return crypto.EcdsaHexToPrivateKey(val)
}

func (e *EnvLoader) LoadRsaPrivateKey(key string, panicOnMissing bool) *crypto.RsaPrivateKey {
	val := os.Getenv(key)
	if val == "" {
		if panicOnMissing {
			panic("missing env var: " + key)
		} else {
			return nil
		}
	}

	return crypto.RsaLoadPrivateKeyFromString(val)
}

func (e *EnvLoader) LoadRsaPublicKey(key string, panicOnMissing bool) *crypto.RsaPublicKey {
	val := os.Getenv(key)
	if val == "" {
		if panicOnMissing {
			panic("missing env var: " + key)
		} else {
			return nil
		}
	}

	return crypto.RsaLoadPublicKeyFromString(val)
}
