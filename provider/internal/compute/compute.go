package compute

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matzapata/ipfs-compute/provider/internal/artifact"
	"github.com/matzapata/ipfs-compute/provider/pkg/escrow"
)

type IComputeService interface {
	Compute(cid string, payerHeader string, computeArgs string) (res *ComputeResponse, ctx *ComputeContext, err error)
	ExecuteProgram(deploymentPath string, execEnv []string, execArgs string) (*ComputeResponse, error)
}

type ComputeService struct {
	ArtifactService         *artifact.ArtifactService
	EscrowService           *escrow.EscrowService
	ProviderEcdsaPrivateKey *ecdsa.PrivateKey
	ProviderEcdsaAddress    *common.Address
	ProviderRsaPrivateKey   *rsa.PrivateKey
	PriceUnit               *big.Int
}

type ComputeResponse struct {
	Data    string            `json:"data"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
}

type ComputeContext struct {
	EscrowTransaction string `json:"escrow_transaction"`
}

func NewComputeService(
	// Services
	artifactService *artifact.ArtifactService,
	escrowService *escrow.EscrowService,

	// Config
	providerEcdsaPrivateKey *ecdsa.PrivateKey,
	providerEcdsaAddress *common.Address,
	providerRsaPrivateKey *rsa.PrivateKey,
	priceUnit *big.Int,
) *ComputeService {
	return &ComputeService{
		ArtifactService:         artifactService,
		EscrowService:           escrowService,
		ProviderRsaPrivateKey:   providerRsaPrivateKey,
		ProviderEcdsaAddress:    providerEcdsaAddress,
		ProviderEcdsaPrivateKey: providerEcdsaPrivateKey,
		PriceUnit:               priceUnit,
	}
}

func (c *ComputeService) Compute(cid string, payerHeader string, computeArgs string) (res *ComputeResponse, ctx *ComputeContext, err error) {
	artifact, err := c.ArtifactService.GetArtifactSpecification(cid, c.ProviderRsaPrivateKey)
	if err != nil {
		return
	}

	// extract payer from header or deployment to owner
	var payerAddress common.Address
	if payerHeader != "" {
		// TODO: extract signer from signature and verify it (Signature has expiration time)
		panic("Not implemented")
	} else {
		payerAddress = common.HexToAddress(artifact.Owner)
	}

	// check if deployment can be paid and other prechecks before executing the binary
	allowance, price, err := c.EscrowService.Allowance(payerAddress, *c.ProviderEcdsaAddress)
	if err != nil {
		return
	}
	if allowance.Cmp(price) < 0 {
		err = errors.New("insufficient funds")
		return
	}
	if price.Cmp(c.PriceUnit) > 0 {
		err = errors.New("invalid price")
		return
	}

	// get deployment executable
	executableDir, err := c.ArtifactService.GetArtifactExecutable(artifact.DeploymentCid)
	if err != nil {
		return
	}

	// reduce user balance in escrow contract.
	tx, err := c.EscrowService.Consume(c.ProviderEcdsaPrivateKey, payerAddress, c.PriceUnit)
	if err != nil {
		return
	}

	// execute binary and give response
	res, err = c.ExecuteProgram(executableDir, artifact.Env, computeArgs)
	if err != nil {
		return
	}
	ctx = &ComputeContext{EscrowTransaction: tx}

	return res, ctx, nil

}

func (c *ComputeService) ExecuteProgram(deploymentPath string, execEnv []string, execArgs string) (*ComputeResponse, error) {
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

	var response ComputeResponse
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %v", err)
	}

	return &response, nil
}
