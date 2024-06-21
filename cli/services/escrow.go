package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/contracts"
)

type EscrowService struct {
	Escrow        *contracts.Escrow
	EthClient     *ethclient.Client
	EscrowAddress string
	Usdc          *contracts.Erc
}

func NewEscrowService(client *ethclient.Client, escrowAddress string, usdcAddress string) *EscrowService {
	address := common.HexToAddress(escrowAddress)
	escrow, err := contracts.NewEscrow(address, client)
	if err != nil {
		panic(err)
	}

	Usdc, err := contracts.NewErc(common.HexToAddress(usdcAddress), client)
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

func (s *EscrowService) ApproveEscrow(hexPrivateKey string, amount uint) (string, error) {
	auth, err := buildAuth(hexPrivateKey, s.EthClient)
	if err != nil {
		return "", err
	}

	// Approve escrow contract to spend USDC
	approveTx, err := s.Usdc.Approve(auth, common.HexToAddress(s.EscrowAddress), big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return approveTx.Hash().Hex(), nil
}

func (s *EscrowService) Deposit(hexPrivateKey string, amount uint) (string, error) {
	// check if the user has approved the escrow contract to spend USDC
	allowance, err := s.Usdc.Allowance(nil, common.HexToAddress(s.EscrowAddress), common.HexToAddress(s.EscrowAddress))
	if err != nil {
		return "", err
	}
	if allowance.Cmp(big.NewInt(int64(amount))) < 0 {
		return "", errors.New("escrow contract is not approved to spend USDC")
	}

	// Deposit funds into the escrow account
	auth, err := buildAuth(hexPrivateKey, s.EthClient)
	if err != nil {
		return "", err
	}
	depositTx, err := s.Escrow.Deposit(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return depositTx.Hash().Hex(), err
}

func (s *EscrowService) Withdraw(hexPrivateKey string, amount uint) (string, error) {
	auth, err := buildAuth(hexPrivateKey, s.EthClient)
	if err != nil {
		return "", err
	}

	tx, err := s.Escrow.Withdraw(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), err
}

func (s *EscrowService) Allowance(userAddress string, providerAddress string) (*big.Int, *big.Int, error) {
	userAccount := common.HexToAddress(userAddress)
	providerAccount := common.HexToAddress(providerAddress)

	return s.Escrow.Allowance(nil, userAccount, providerAccount)
}

func (s *EscrowService) Balance(userAddress string) (*big.Int, error) {
	account := common.HexToAddress(userAddress)

	return s.Escrow.BalanceOf(nil, account)
}

func (s *EscrowService) ApproveProvider(hexPrivateKey string, providerAddress string, amount uint, price uint) (string, error) {
	auth, err := buildAuth(hexPrivateKey, s.EthClient)
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

func buildAuth(hexPrivateKey string, client *ethclient.Client) (*bind.TransactOpts, error) {
	// recover the public key from the private key
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	ownerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Some preparation for the transaction
	nonce, err := client.PendingNonceAt(context.Background(), ownerAddress)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	chainId, _ := client.ChainID(context.Background())
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice

	return auth, nil
}
