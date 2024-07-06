package domain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type DepositTransactions struct {
	Approval string
	Deposit  string
}

type IEscrowService interface {
	Consume(privateKey *crypto.EcdsaPrivateKey, userAddress crypto.EcdsaAddress, priceUnit *big.Int) (string, error)
	Deposit(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (*DepositTransactions, error)
	Withdraw(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (string, error)
	Allowance(userAddress crypto.EcdsaAddress, providerAddress crypto.EcdsaAddress) (*big.Int, *big.Int, error)
	Balance(userAddress crypto.EcdsaAddress) (*big.Int, error)
	ApproveProvider(privateKey *crypto.EcdsaPrivateKey, providerAddress crypto.EcdsaAddress, amount *big.Int, price *big.Int) (string, error)
}

type IEscrowContract interface {
	Consume(opts *bind.TransactOpts, account crypto.EcdsaAddress, price *big.Int) (*types.Transaction, error)
	Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)
	Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)
	Allowance(opts *bind.CallOpts, user crypto.EcdsaAddress, provider crypto.EcdsaAddress) (*big.Int, *big.Int, error)
	BalanceOf(opts *bind.CallOpts, account crypto.EcdsaAddress) (*big.Int, error)
	Approve(opts *bind.TransactOpts, provider crypto.EcdsaAddress, amount *big.Int, price *big.Int) (*types.Transaction, error)
}
