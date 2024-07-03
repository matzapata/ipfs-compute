package domain

import (
	"math/big"

	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type IEscrowContract interface {
}

type IEscrowService interface {
	Consume(privateKey *crypto.EcdsaPrivateKey, userAddress crypto.EcdsaAddress, priceUnit *big.Int) (string, error)
	ApproveEscrow(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (string, error)
	Deposit(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (string, error)
	Withdraw(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (string, error)
	Allowance(userAddress crypto.EcdsaAddress, providerAddress crypto.EcdsaAddress) (*big.Int, *big.Int, error)
	Balance(userAddress string) (*big.Int, error)
	ApproveProvider(privateKey *crypto.EcdsaPrivateKey, providerAddress crypto.EcdsaAddress, amount *big.Int, price *big.Int) (string, error)
}
