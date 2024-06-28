package services

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/internal/contracts"
	eth_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/eth"
)

type EscrowService struct {
	Escrow        *contracts.Escrow
	EthClient     *ethclient.Client
	EscrowAddress *common.Address
	Usdc          *contracts.Erc
}

func NewEscrowService(client *ethclient.Client, escrowAddress *common.Address, usdcAddress *common.Address) *EscrowService {
	escrow, err := contracts.NewEscrow(*escrowAddress, client)
	if err != nil {
		panic(err)
	}

	Usdc, err := contracts.NewErc(*usdcAddress, client)
	if err != nil {
		panic(err)
	}

	return &EscrowService{
		Escrow:        escrow,
		EthClient:     client,
		EscrowAddress: escrowAddress,
		Usdc:          Usdc,
	}
}

func (s *EscrowService) Consume(privateKey *ecdsa.PrivateKey, userAddress common.Address, priceUnit *big.Int) (string, error) {
	auth, err := eth_helpers.BuildAuth(privateKey, s.EthClient, nil)
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

func (s *EscrowService) ApproveEscrow(privateKey *ecdsa.PrivateKey, amount uint) (string, error) {
	auth, err := eth_helpers.BuildAuth(privateKey, s.EthClient, nil)
	if err != nil {
		return "", err
	}

	// Approve escrow contract to spend USDC
	approveTx, err := s.Usdc.Approve(auth, *s.EscrowAddress, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return approveTx.Hash().Hex(), nil
}

func (s *EscrowService) Deposit(privateKey *ecdsa.PrivateKey, amount uint) (string, error) {
	// extract owner from private key
	ownerAddress, err := eth_helpers.PrivateKeyToAddress(privateKey)
	if err != nil {
		return "", err
	}

	// check if the user has approved the escrow contract to spend USDC
	allowance, err := s.Usdc.Allowance(nil, ownerAddress, *s.EscrowAddress)
	if err != nil {
		return "", err
	}
	if allowance.Cmp(big.NewInt(int64(amount))) < 0 {
		return "", errors.New("escrow contract is not approved to spend USDC")
	}

	// Deposit funds into the escrow account
	auth, err := eth_helpers.BuildAuth(privateKey, s.EthClient, nil)
	if err != nil {
		return "", err
	}
	depositTx, err := s.Escrow.Deposit(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return depositTx.Hash().Hex(), err
}

func (s *EscrowService) Withdraw(privateKey *ecdsa.PrivateKey, amount uint) (string, error) {
	auth, err := eth_helpers.BuildAuth(privateKey, s.EthClient, nil)
	if err != nil {
		return "", err
	}

	tx, err := s.Escrow.Withdraw(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), err
}

func (s *EscrowService) Allowance(userAddress common.Address, providerAddress common.Address) (*big.Int, *big.Int, error) {
	return s.Escrow.Allowance(nil, userAddress, providerAddress)
}

func (s *EscrowService) Balance(userAddress string) (*big.Int, error) {
	account := common.HexToAddress(userAddress)

	return s.Escrow.BalanceOf(nil, account)
}

func (s *EscrowService) ApproveProvider(privateKey *ecdsa.PrivateKey, providerAddress string, amount uint, price uint) (string, error) {
	auth, err := eth_helpers.BuildAuth(privateKey, s.EthClient, nil)
	if err != nil {
		return "", err
	}

	// Approve provider to spend USDC
	tx, err := s.Escrow.Approve(auth, common.HexToAddress(providerAddress), big.NewInt(int64(amount)), big.NewInt(int64(price)))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
