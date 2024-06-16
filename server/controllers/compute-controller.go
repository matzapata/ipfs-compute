package controllers

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/helpers"
	"github.com/matzapata/ipfs-compute/services"
)

type ComputeHttpController struct {
	computeService    *services.ComputeService
	deploymentService *services.DeploymentService
}

func NewComputeHttpController(computeService *services.ComputeService, deploymentService *services.DeploymentService) *ComputeHttpController {
	return &ComputeHttpController{
		computeService:    computeService,
		deploymentService: deploymentService,
	}
}

func (c *ComputeHttpController) Compute(w http.ResponseWriter, r *http.Request) {
	// extract data from request
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get deployment specifications
	depl, err := c.deploymentService.GetDeploymentMetadata(cid)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// TODO: verify the signature of the deployment

	// TODO: check if deployment can be paid and other prechecks before executing the binary

	// temp folder creation
	tempDir, err := os.MkdirTemp("", "ipfs-compute-*")
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// download deployment
	err = c.deploymentService.GetDeployment(depl.DeploymentCid, tempDir)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// TODO: reduce user balance in escrow contract. Extract payer from header or deployment

	// execute binary and give response
	resultJSON, err := c.computeService.Compute(tempDir, depl.Env)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}
	helpers.WriteJSON(w, resultJSON.Status, resultJSON.Data, resultJSON.Headers)
}
