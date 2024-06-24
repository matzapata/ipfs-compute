package compute

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matzapata/ipfs-compute/provider/internal/deployment"
	"github.com/matzapata/ipfs-compute/provider/pkg/escrow"
)

type ComputeService struct {
	DeploymentService       *deployment.DeploymentService
	EscrowService           *escrow.EscrowService
	ProviderEcdsaAddress    *common.Address
	ProviderEcdsaPrivateKey *ecdsa.PrivateKey
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
	deploymentService *deployment.DeploymentService,
	escrowService *escrow.EscrowService,
	providerEcdsaAddress *common.Address,
	providerEcdsaPrivateKey *ecdsa.PrivateKey,
	priceUnit *big.Int,
) *ComputeService {
	return &ComputeService{
		DeploymentService:    deploymentService,
		EscrowService:        escrowService,
		ProviderEcdsaAddress: providerEcdsaAddress,
		PriceUnit:            priceUnit,
	}
}

func (c *ComputeService) Compute(cid string, payerHeader string, computeArgs string) (*ComputeResponse, *ComputeContext, error) {
	// Get deployment specifications
	depl, err := c.DeploymentService.GetDeploymentMetadata(cid)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Extract payer from header or deployment to owner
	var payerAddress common.Address
	if payerHeader != "" {
		// TODO: extract signer from signature and verify it (Signature has expiration time)
		panic("Not implemented")
	} else {
		payerAddress = common.HexToAddress(depl.Owner)
	}

	// check if deployment can be paid and other prechecks before executing the binary
	allowance, price, err := c.EscrowService.Allowance(payerAddress, *c.ProviderEcdsaAddress)
	if err != nil {
		return nil, nil, err
	}
	if allowance.Cmp(price) < 0 {
		return nil, nil, errors.New("insufficient funds")
	}
	if price.Cmp(c.PriceUnit) > 0 {
		return nil, nil, errors.New("invalid price")
	}

	// temp folder creation.
	tempDir, err := os.MkdirTemp("", "khachapuri-*")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(tempDir)

	// download deployment
	err = c.DeploymentService.GetDeployment(depl.DeploymentCid, tempDir)
	if err != nil {
		return nil, nil, err
	}

	// reduce user balance in escrow contract.
	tx, err := c.EscrowService.Consume(c.ProviderEcdsaPrivateKey, payerAddress, c.PriceUnit)
	if err != nil {
		return nil, nil, err
	}

	// execute binary and give response
	resultJSON, err := c.ExecuteProgram(tempDir, depl.Env, computeArgs)
	if err != nil {
		return nil, nil, err
	}

	context := &ComputeContext{
		EscrowTransaction: tx,
	}

	return resultJSON, context, nil

}

func (cs *ComputeService) ExecuteProgram(deploymentPath string, execEnv []string, execArgs string) (*ComputeResponse, error) {
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

// Creates curl like command to be executed in the gateway
func (cs *ComputeService) ParseRequest(r *http.Request) (string, error) {
	// Prepare the curl command
	args := []string{"-X", r.Method}

	// add headers
	for key, value := range r.Header {
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value[0]))
	}
	args = append(args, r.URL.String())

	// add data
	if r.Method == "POST" || r.Method == "PUT" {
		body, err := r.GetBody()
		if err != nil {
			return "", fmt.Errorf("failed to get body: %v", err)
		}
		data, err := io.ReadAll(body)
		if err != nil {
			return "", fmt.Errorf("failed to read body: %v", err)
		}
		args = append(args, "-d", string(data))
	}

	// TODO: add more methods

	return fmt.Sprintf("%s", args), nil
}
