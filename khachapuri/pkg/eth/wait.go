package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITransactionWaiter interface {
	Wait(hash common.Hash) (*types.Receipt, error)
}

type TransactionWaiter struct {
	EthClient *ethclient.Client
}

func NewTransactionWaiter(ethClient *ethclient.Client) *TransactionWaiter {
	return &TransactionWaiter{
		EthClient: ethClient,
	}
}

func (t *TransactionWaiter) Wait(hash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	for {
		receipt, err := t.EthClient.TransactionReceipt(ctx, hash)
		if err == ethereum.NotFound {
			// Transaction not yet mined, wait and try again
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			return nil, err
		}
		return receipt, nil
	}
}
