package domain

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type IEscrowService interface {
	Consume(privateKey *ecdsa.PrivateKey, userAddress common.Address, priceUnit *big.Int) (string, error)
	ApproveEscrow(privateKey *ecdsa.PrivateKey, amount uint) (string, error)
	Deposit(privateKey *ecdsa.PrivateKey, amount uint) (string, error)
	Withdraw(privateKey *ecdsa.PrivateKey, amount uint) (string, error)
	Allowance(userAddress common.Address, providerAddress common.Address) (*big.Int, *big.Int, error)
	Balance(userAddress string) (*big.Int, error)
	ApproveProvider(privateKey *ecdsa.PrivateKey, providerAddress string, amount uint, price uint) (string, error)
}
