package api_routers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	api_helpers "github.com/matzapata/ipfs-compute/provider/internal/controllers/api/helpers"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
)

func SetupComputeRoutes(router *chi.Mux, computeService domain.IComputeService) {
	router.HandleFunc("/{cid}", func(w http.ResponseWriter, r *http.Request) {
		cid := chi.URLParam(r, "cid")
		if cid == "" {
			api_helpers.ErrorJSON(w, errors.New("CID is required"), http.StatusBadRequest)
			return
		}
		payerHeader := r.Header.Get("x-payer-signature")

		// parse request to executable args. essentially convert request to curl command
		args, err := ParseRequest(r)
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

func ParseRequest(r *http.Request) (string, error) {
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
