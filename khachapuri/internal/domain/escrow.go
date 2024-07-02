package domain

import (
	"math/big"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type IEscrowService interface {
	Consume(privateKey *crypto.EcdsaPrivateKey, userAddress crypto.EcdsaAddress, priceUnit *big.Int) (string, error)
	ApproveEscrow(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error)
	Deposit(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error)
	Withdraw(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error)
	Allowance(userAddress crypto.EcdsaAddress, providerAddress crypto.EcdsaAddress) (*big.Int, *big.Int, error)
	Balance(userAddress string) (*big.Int, error)
	ApproveProvider(privateKey *crypto.EcdsaPrivateKey, providerAddress crypto.EcdsaAddress, amount uint, price uint) (string, error)
}
