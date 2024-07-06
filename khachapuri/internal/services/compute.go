package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/matzapata/ipfs-compute/provider/pkg/crypto"
)

type ComputeService struct {
	ArtifactService         domain.IArtifactService
	EscrowService           domain.IEscrowService
	ProviderEcdsaPrivateKey *crypto.EcdsaPrivateKey
	ProviderEcdsaAddress    *crypto.EcdsaAddress
	ProviderRsaPrivateKey   *crypto.RsaPrivateKey
	ProviderRsaPublicKey    *crypto.RsaPublicKey
	PriceUnit               *big.Int
}

func NewComputeService(
	// Services
	artifactService domain.IArtifactService,
	escrowService domain.IEscrowService,

	// Config
	providerEcdsaPrivateKey *crypto.EcdsaPrivateKey,
	providerEcdsaAddress *crypto.EcdsaAddress,
	providerRsaPrivateKey *crypto.RsaPrivateKey,
	providerRsaPublicKey *crypto.RsaPublicKey,
	priceUnit *big.Int,
) *ComputeService {
	return &ComputeService{
		ArtifactService:         artifactService,
		EscrowService:           escrowService,
		ProviderRsaPrivateKey:   providerRsaPrivateKey,
		ProviderRsaPublicKey:    providerRsaPublicKey,
		ProviderEcdsaPrivateKey: providerEcdsaPrivateKey,
		ProviderEcdsaAddress:    providerEcdsaAddress,
		PriceUnit:               priceUnit,
	}
}

func (c *ComputeService) Compute(cid string, payerHexAddress string, computeArgs string) (res *domain.ComputeResponse, ctx *domain.ComputeContext, err error) {
	// download specification
	artifact, err := c.ArtifactService.GetArtifactSpecification(cid, c.ProviderRsaPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	// by default payer is owner
	var payer crypto.EcdsaAddress
	if payerHexAddress != "" {
		payer = common.HexToAddress(payerHexAddress)
	} else {
		payer = common.HexToAddress(artifact.Owner)
	}

	// check if deployment can be paid and other prechecks before executing the binary
	allowance, price, err := c.EscrowService.Allowance(payer, *c.ProviderEcdsaAddress)
	if err != nil {
		return nil, nil, err
	}
	if allowance.Cmp(price) < 0 {
		return nil, nil, errors.New("insufficient funds")
	}
	if price.Cmp(c.PriceUnit) > 0 {
		return nil, nil, errors.New("invalid price")
	}

	// get deployment executable
	executableDir, err := c.ArtifactService.GetArtifactExecutable(artifact.DeploymentCid)
	if err != nil {
		return nil, nil, err
	}

	// reduce user balance in escrow contract.
	tx, err := c.EscrowService.Consume(c.ProviderEcdsaPrivateKey, payer, c.PriceUnit)
	if err != nil {
		return nil, nil, err
	}

	// execute binary and give response
	res, err = ExecuteProgram(executableDir, artifact.Env, computeArgs)
	if err != nil {
		return nil, nil, err
	}
	ctx = &domain.ComputeContext{EscrowTransaction: tx}

	return res, ctx, nil

}

func ExecuteProgram(deploymentPath string, execEnv []string, execArgs string) (*domain.ComputeResponse, error) {
	// Prepare the docker run command
	args := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/app", deploymentPath)}
	for _, env := range execEnv {
		args = append(args, "-e", env)
	}
	args = append(args, "binary_runner", "main") // binary_runner is the name of the docker image
	args = append(args, execArgs)

	// run the binary inside the docker container
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("execution error: %v, stderr: %v", err, stderr.String())
	}

	var response domain.ComputeResponse
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %v", err)
	}

	return &response, nil
}
