package domain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type IErcAllowance interface {
	Approve(opts *bind.TransactOpts, spender crypto.EcdsaAddress, value *big.Int) (*types.Transaction, error)
	Allowance(opts *bind.CallOpts, owner crypto.EcdsaAddress, spender crypto.EcdsaAddress) (*big.Int, error)
}
