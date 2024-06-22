package services

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/server/config"
	"github.com/matzapata/ipfs-compute/shared/contracts"
)

type EscrowService struct {
	Escrow        *contracts.Escrow
	EthClient     *ethclient.Client
	EscrowAddress string
}

func NewEscrowService(client *ethclient.Client, escrowAddress string) *EscrowService {
	address := common.HexToAddress(escrowAddress)
	escrow, err := contracts.NewEscrow(address, client)
	if err != nil {
		panic(err)
	}

	return &EscrowService{
		Escrow:        escrow,
		EthClient:     client,
		EscrowAddress: escrowAddress,
	}
}

func (s *EscrowService) Consume(hexPrivateKey string, userAddress common.Address, priceUnit *big.Int) (string, error) {
	auth, err := buildAuth(hexPrivateKey, s.EthClient)
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

func (s *EscrowService) Allowance(address common.Address) (*big.Int, *big.Int, error) {
	providerAccount := common.HexToAddress("providerAddress")

	return s.Escrow.Allowance(nil, address, providerAccount)

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
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, config.CHAIN_ID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice

	return auth, nil
}
