package api_routers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/compute"
	api_helpers "github.com/matzapata/ipfs-compute/provider/internal/controllers/api/helpers"
)

func SetupComputeRoutes(router *chi.Mux, computeService compute.ComputeService) {
	router.HandleFunc("/{cid}", func(w http.ResponseWriter, r *http.Request) {
		cid := chi.URLParam(r, "cid")
		if cid == "" {
			api_helpers.ErrorJSON(w, errors.New("CID is required"), http.StatusBadRequest)
			return
		}
		payerHeader := r.Header.Get("x-payer-signature")

		// parse request to executable args. essentially convert request to curl command
		args, err := computeService.ParseRequest(r)
		if err != nil {
			api_helpers.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// execute compute
		res, ctx, err := computeService.Compute(cid, payerHeader, args)
		if err != nil {
			api_helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		// Add execution context to headers
		res.Headers["x-escrow-tx"] = ctx.EscrowTransaction

		// write response
		api_helpers.WriteJSON(w, res.Status, res.Data, res.Headers)
	})
}
