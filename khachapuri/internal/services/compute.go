package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type ComputeService struct {
	Config          *config.Config
	ArtifactService domain.IArtifactService
	EscrowService   domain.IEscrowService
	ComputeExecutor domain.IComputeExecutor
}

func NewComputeService(
	cfg *config.Config,
	artifactService domain.IArtifactService,
	escrowService domain.IEscrowService,
	computeExecutor domain.IComputeExecutor,
) *ComputeService {
	return &ComputeService{
		ArtifactService: artifactService,
		EscrowService:   escrowService,
		Config:          cfg,
		ComputeExecutor: computeExecutor,
	}
}

func (c *ComputeService) Compute(specCid string, payerHexAddress string, computeArgs string) (res *domain.ComputeResponse, ctx *domain.ComputeContext, err error) {
	defer func() {
		// recover from panic to properly return ctx
		if r := recover(); r != nil {
			err = fmt.Errorf("error: %v", r)
		}
	}()

	// download specification
	artifact, err := c.ArtifactService.GetArtifactSpecification(specCid, c.Config.ProviderRsaPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	execEnv := []string{}
	execEnv = append(execEnv, strings.Split(strings.TrimSpace(string(artifact.Env)), "\n")...)

	// by default payer is owner
	var payer crypto.EcdsaAddress
	if payerHexAddress != "" {
		payer = common.HexToAddress(payerHexAddress)
	} else {
		payer = common.HexToAddress(artifact.Owner)
	}

	// check if deployment can be paid and other pre-checks before executing the binary
	allowance, price, err := c.EscrowService.Allowance(payer, *c.Config.ProviderEcdsaAddress)
	if err != nil {
		return nil, nil, err
	}
	if allowance.Cmp(price) < 0 {
		return nil, nil, errors.New("insufficient funds")
	}
	if price.Cmp(c.Config.ProviderComputeUnitPrice) > 0 {
		return nil, nil, errors.New("invalid price")
	}

	// get deployment executable
	executableDir, err := c.ArtifactService.GetArtifact(artifact.ArtifactCid)
	if err != nil {
		return nil, nil, err
	}

	// reduce user balance in escrow contract.
	tx, err := c.EscrowService.Consume(c.Config.ProviderEcdsaPrivateKey, payer, c.Config.ProviderComputeUnitPrice)
	if err != nil {
		return nil, nil, err
	}
	ctx = &domain.ComputeContext{EscrowTransaction: tx}

	// execute binary and give response
	res, err = c.ComputeExecutor.Execute(executableDir, execEnv, computeArgs)
	if err != nil {
		return nil, ctx, err
	}

	return res, ctx, nil

}
