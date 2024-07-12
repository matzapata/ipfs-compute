package routers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
)

func SetupComputeRoutes(router *chi.Mux, computeService domain.IComputeService) {
	router.HandleFunc("/{cid}/*", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, computeService)
	})
	router.HandleFunc("/{cid}", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, computeService)
	})
}

func ParseRequest(r *http.Request, remainingPath string) (string, error) {
	// Prepare the curl command
	args := []string{"-X", r.Method}

	// add headers
	for key, value := range r.Header {
		args = append(args, "-H", fmt.Sprintf("\"%s: %s\"", key, value[0]))
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

	args = append(args, remainingPath)

	return fmt.Sprintf("%s", args), nil
}

func handler(w http.ResponseWriter, r *http.Request, computeService domain.IComputeService) {
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		ErrorJSON(w, http.StatusBadRequest, errors.New("CID is required"), nil)
		return
	}
	fmt.Println(r.URL.Path, cid)
	remainingPath := strings.TrimPrefix(r.URL.Path, "/"+cid)
	payerHeader := r.Header.Get("x-payer-signature")

	// parse request to executable args. essentially convert request to curl command
	args, err := ParseRequest(r, remainingPath)
	if err != nil {
		ErrorJSON(w, http.StatusBadRequest, err, nil)
		return
	}

	// execute compute
	res, ctx, err := computeService.Compute(cid, payerHeader, args)

	// assemble headers
	headers := map[string]string{
		"x-escrow-tx": ctx.EscrowTransaction,
	}
	if res != nil {
		for key, value := range res.Headers {
			headers[key] = value
		}
	}

	fmt.Println(ctx.EscrowTransaction)
	if err != nil {
		ErrorJSON(w, http.StatusInternalServerError, err, headers)
	} else {
		WriteJSON(w, res.Status, res.Data, headers)
	}
}
