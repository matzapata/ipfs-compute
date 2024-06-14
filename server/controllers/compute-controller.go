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

	// temp folder creation
	tempDir, err := os.MkdirTemp("", "ipfs-compute-*")
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// get deployment
	depl, err := c.deploymentService.GetDeployment(cid, tempDir)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// reduce user balance in escrow contract. Extract payer from header or deployment

	// execute binary and give response
	resultJSON, err := c.computeService.Compute(tempDir, depl.Env)
	if err != nil {
		helpers.ErrorJSON(w, err, 400)
		return
	}
	helpers.WriteJSON(w, resultJSON.Status, resultJSON.Data, resultJSON.Headers)
}
