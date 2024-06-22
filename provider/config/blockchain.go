package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var CHAIN_ID = new(big.Int).SetInt64(137)

var ESCROW_ADDRESS = common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

var USDC_ADDRESS = common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
