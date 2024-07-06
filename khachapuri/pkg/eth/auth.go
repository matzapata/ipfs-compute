package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	goeth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type IEthAuthenticator interface {
	Authenticate(privateKey *crypto.EcdsaPrivateKey) (*bind.TransactOpts, error)
}

type EthAuthenticator struct {
	EthClient *ethclient.Client
}

func NewEthAuthenticator(ethClient *ethclient.Client) *EthAuthenticator {
	return &EthAuthenticator{
		EthClient: ethClient,
	}
}

func (ea *EthAuthenticator) Authenticate(privateKey *crypto.EcdsaPrivateKey) (*bind.TransactOpts, error) {
	// recover the public key from the private key
	ownerAddress, err := PrivateKeyToAddress(privateKey)
	if err != nil {
		return nil, err
	}

	// Some preparation for the transaction
	nonce, err := ea.EthClient.PendingNonceAt(context.Background(), ownerAddress)
	if err != nil {
		return nil, err
	}
	gasPrice, err := ea.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	chainId, err := ea.EthClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice

	return auth, nil
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	ownerAddress := goeth_crypto.PubkeyToAddress(*publicKeyECDSA)
	return ownerAddress, nil
}
