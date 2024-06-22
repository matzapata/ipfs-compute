package controllers

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/services"
	"github.com/matzapata/ipfs-compute/provider/utils"
	"github.com/matzapata/ipfs-compute/shared/cryptoecdsa"
	"github.com/matzapata/ipfs-compute/shared/escrow"
)

type ComputeHttpController struct {
	ComputeService          *services.ComputeService
	DeploymentService       *services.DeploymentService
	EscrowService           *escrow.EscrowService
	PriceUnit               *big.Int
	ProviderEcdsaPrivateKey *ecdsa.PrivateKey
	ProviderEcdsaAddress    *common.Address
}

func NewComputeHttpController(
	computeService *services.ComputeService,
	deploymentService *services.DeploymentService,
	escrowService *escrow.EscrowService,
	priceUnit *big.Int,
	providerEcdsaPrivateKey *ecdsa.PrivateKey,
) *ComputeHttpController {
	providerAddress, err := cryptoecdsa.PrivateKeyToAddress(providerEcdsaPrivateKey)
	if err != nil {
		panic(err)
	}

	return &ComputeHttpController{
		ComputeService:          computeService,
		DeploymentService:       deploymentService,
		EscrowService:           escrowService,
		PriceUnit:               priceUnit,
		ProviderEcdsaPrivateKey: providerEcdsaPrivateKey,
		ProviderEcdsaAddress:    &providerAddress,
	}
}

func (c *ComputeHttpController) Compute(w http.ResponseWriter, r *http.Request) {
	// extract data from request
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		utils.ErrorJSON(w, errors.New("CID is required"), http.StatusBadRequest)
		return
	}
	payerHeader := r.Header.Get("x-payer-signature")

	// Get deployment specifications
	depl, err := c.DeploymentService.GetDeploymentMetadata(cid)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
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
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if allowance.Cmp(price) < 0 {
		utils.ErrorJSON(w, errors.New("insufficient funds"), http.StatusPaymentRequired)
		return
	}
	if price.Cmp(c.PriceUnit) > 0 {
		utils.ErrorJSON(w, errors.New("invalid price"), http.StatusBadRequest)
		return
	}

	// temp folder creation.
	tempDir, err := os.MkdirTemp("", "khachapuri-*")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// download deployment
	err = c.DeploymentService.GetDeployment(depl.DeploymentCid, tempDir)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// reduce user balance in escrow contract.
	tx, err := c.EscrowService.Consume(c.ProviderEcdsaPrivateKey, payerAddress, c.PriceUnit)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// execute binary and give response
	resultJSON, err := c.ComputeService.Compute(tempDir, depl.Env)
	if err != nil {
		utils.ErrorJSON(w, err, 400)
		return
	}

	// Add execution context to headers
	resultJSON.Headers["x-escrow-tx"] = tx

	// write response
	utils.WriteJSON(w, resultJSON.Status, resultJSON.Data, resultJSON.Headers)
}
