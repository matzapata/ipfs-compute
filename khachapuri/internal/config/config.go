package config

import (
	"math/big"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type Config struct {
	EthRpc                   string
	IpfsGateway              string
	IpfsPinataApikey         string
	IpfsPinataSecret         string
	ArtifactMaxSize          uint
	ArtifactsPath            string
	CachePath                string
	TempPath                 string
	ProviderEcdsaAddress     *crypto.EcdsaAddress
	ProviderEcdsaPrivateKey  *crypto.EcdsaPrivateKey
	ProviderRsaPrivateKey    *crypto.RsaPrivateKey
	ProviderRsaPublicKey     *crypto.RsaPublicKey
	ProviderComputeUnitPrice *big.Int
	EscrowAddress            *crypto.EcdsaAddress
	UsdcAddress              *crypto.EcdsaAddress
	RegistryAddress          *crypto.EcdsaAddress
}

type ConfigLoader interface {
	LoadString(key string, panicOnMissing bool) string
	LoadInt(key string, panicOnMissing bool) int
	LoadBool(key string, panicOnMissing bool) bool
	LoadBigInt(key string, panicOnMissing bool) *big.Int
	LoadEcdsaAddress(key string, panicOnMissing bool) *crypto.EcdsaAddress
	LoadEcdsaPrivateKey(key string, panicOnMissing bool) *crypto.EcdsaPrivateKey
	LoadRsaPrivateKey(key string, panicOnMissing bool) *crypto.RsaPrivateKey
	LoadRsaPublicKey(key string, panicOnMissing bool) *crypto.RsaPublicKey
}

var UsdcAddress = crypto.EcdsaHexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
var RegistryAddress = crypto.EcdsaHexToAddress("0xdb42A86B1bfe04E75B2A5F2bF7a3BBB52D7FFD2F")
var EscrowAddress = crypto.EcdsaHexToAddress("0x5Fe8861F6571174a9564365384AE9b01CcdCd8D6")
var ArtifactMaxSize uint = 50 * 1024 * 1024
