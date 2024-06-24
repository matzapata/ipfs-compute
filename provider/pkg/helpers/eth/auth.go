package eth_helpers

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	ownerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return ownerAddress, nil
}

func BuildAuth(privateKey *ecdsa.PrivateKey, client *ethclient.Client, chainID *big.Int) (*bind.TransactOpts, error) {
	// recover the public key from the private key
	ownerAddress, err := PrivateKeyToAddress(privateKey)
	if err != nil {
		return nil, err
	}

	// Some preparation for the transaction
	nonce, err := client.PendingNonceAt(context.Background(), ownerAddress)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	var chainId *big.Int
	if chainID == nil {
		chainId, err = client.ChainID(context.Background())
		if err != nil {
			return nil, err
		}
	} else {
		chainId = chainID
	}

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice

	return auth, nil
}
