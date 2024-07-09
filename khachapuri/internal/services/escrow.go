package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type EscrowService struct {
	Config            *config.Config
	Escrow            domain.IEscrowContract
	Erc               domain.IErcAllowance
	Authenticator     eth.IEthAuthenticator
	TransactionWaiter eth.ITransactionWaiter
}

func NewEscrowService(
	cfg *config.Config,
	ethClient *ethclient.Client,
) *EscrowService {
	ercContract, err := contracts.NewErc(*cfg.UsdcAddress, ethClient)
	if err != nil {
		panic(err)
	}
	escrowContract, err := contracts.NewEscrow(*cfg.EscrowAddress, ethClient)
	if err != nil {
		panic(err)
	}
	ethAuthenticator := eth.NewEthAuthenticator(ethClient)
	transactionWaiter := eth.NewTransactionWaiter(ethClient)

	return &EscrowService{
		Config:            cfg,
		Escrow:            escrowContract,
		Erc:               ercContract,
		Authenticator:     ethAuthenticator,
		TransactionWaiter: transactionWaiter,
	}
}

// executes transaction in escrow that extracts 1 credit at priceUnit from userAddress
func (s *EscrowService) Consume(providerPrivateKey *crypto.EcdsaPrivateKey, userAddress crypto.EcdsaAddress, priceUnit *big.Int) (string, error) {
	auth, err := s.Authenticator.Authenticate(providerPrivateKey)
	if err != nil {
		return "", err
	}

	// Consume credit from the user
	tx, err := s.Escrow.Consume(auth, userAddress, priceUnit)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil

}

func (s *EscrowService) Deposit(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (*domain.DepositTransactions, error) {
	// recover address from private key
	address := crypto.EcdsaPrivateKeyToAddress(privateKey)

	// check if the user has approved the escrow contract to spend USDC
	allowance, err := s.Erc.Allowance(nil, *address, *s.Config.EscrowAddress)
	if err != nil {
		return nil, err
	}
	approveTx := ""
	if allowance.Cmp(amount) < 0 {
		// Approve escrow contract to spend USDC
		auth, err := s.Authenticator.Authenticate(privateKey)
		if err != nil {
			return nil, err
		}

		tx, err := s.Erc.Approve(auth, *s.Config.EscrowAddress, amount)
		if err != nil {
			return nil, err
		}
		approveTx = tx.Hash().Hex()

		// wait for the approval transaction to be mined
		_, err = s.TransactionWaiter.Wait(tx.Hash())
		if err != nil {
			return &domain.DepositTransactions{
				Approval: approveTx,
			}, err
		}
	}

	// Deposit funds into the escrow account
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return &domain.DepositTransactions{
			Approval: approveTx,
		}, err
	}
	depositTx, err := s.Escrow.Deposit(auth, amount)
	if err != nil {
		return &domain.DepositTransactions{
			Approval: approveTx,
		}, err
	}

	return &domain.DepositTransactions{
		Approval: approveTx,
		Deposit:  depositTx.Hash().Hex(),
	}, nil
}

func (s *EscrowService) Withdraw(privateKey *crypto.EcdsaPrivateKey, amount *big.Int) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	tx, err := s.Escrow.Withdraw(auth, amount)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), err
}

func (s *EscrowService) Allowance(userAddress crypto.EcdsaAddress, providerAddress crypto.EcdsaAddress) (*big.Int, *big.Int, error) {
	return s.Escrow.Allowance(nil, userAddress, providerAddress)
}

func (s *EscrowService) Balance(account crypto.EcdsaAddress) (*big.Int, error) {
	return s.Escrow.BalanceOf(nil, account)
}

func (s *EscrowService) ApproveProvider(privateKey *crypto.EcdsaPrivateKey, providerAddress crypto.EcdsaAddress, amount *big.Int, price *big.Int) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	// Approve provider to spend USDC
	tx, err := s.Escrow.Approve(auth, providerAddress, amount, price)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
