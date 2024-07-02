package services

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
	"github.com/matzapata/ipfs-compute/provider/pkg/eth"
)

type EscrowService struct {
	Escrow        *contracts.Escrow
	EthClient     *ethclient.Client
	Usdc          *contracts.Erc
	Authenticator eth.IEthAuthenticator
}

func NewEscrowService(client *ethclient.Client, ethAuthenticator eth.IEthAuthenticator) *EscrowService {
	escrow, err := contracts.NewEscrow(config.ESCROW_ADDRESS, client)
	if err != nil {
		panic(err)
	}

	Usdc, err := contracts.NewErc(config.USDC_ADDRESS, client)
	if err != nil {
		panic(err)
	}

	return &EscrowService{
		Escrow:        escrow,
		EthClient:     client,
		Usdc:          Usdc,
		Authenticator: ethAuthenticator,
	}
}

func (s *EscrowService) Consume(privateKey *crypto.EcdsaPrivateKey, userAddress crypto.EcdsaAddress, priceUnit *big.Int) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
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

func (s *EscrowService) ApproveEscrow(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	// Approve escrow contract to spend USDC
	approveTx, err := s.Usdc.Approve(auth, config.ESCROW_ADDRESS, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return approveTx.Hash().Hex(), nil
}

func (s *EscrowService) Deposit(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error) {
	// recover address from private key
	address, err := crypto.EcdsaPrivateKeyToAddress(privateKey)
	if err != nil {
		return "", err
	}

	// check if the user has approved the escrow contract to spend USDC
	allowance, err := s.Usdc.Allowance(nil, address, config.ESCROW_ADDRESS)
	if err != nil {
		return "", err
	}
	if allowance.Cmp(big.NewInt(int64(amount))) < 0 {
		return "", errors.New("escrow contract is not approved to spend USDC")
	}

	// Deposit funds into the escrow account
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}
	depositTx, err := s.Escrow.Deposit(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return depositTx.Hash().Hex(), err
}

func (s *EscrowService) Withdraw(privateKey *crypto.EcdsaPrivateKey, amount uint) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	tx, err := s.Escrow.Withdraw(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), err
}

func (s *EscrowService) Allowance(userAddress crypto.EcdsaAddress, providerAddress crypto.EcdsaAddress) (*big.Int, *big.Int, error) {
	return s.Escrow.Allowance(nil, userAddress, providerAddress)
}

func (s *EscrowService) Balance(userAddress string) (*big.Int, error) {
	account := common.HexToAddress(userAddress)

	return s.Escrow.BalanceOf(nil, account)
}

func (s *EscrowService) ApproveProvider(privateKey *crypto.EcdsaPrivateKey, providerAddress crypto.EcdsaAddress, amount uint, price uint) (string, error) {
	auth, err := s.Authenticator.Authenticate(privateKey)
	if err != nil {
		return "", err
	}

	// Approve provider to spend USDC
	tx, err := s.Escrow.Approve(auth, providerAddress, big.NewInt(int64(amount)), big.NewInt(int64(price)))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
