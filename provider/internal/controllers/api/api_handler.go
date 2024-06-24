package api_controller

import (
	"net/http"

	"github.com/go-chi/chi"
	api_routers "github.com/matzapata/ipfs-compute/provider/internal/controllers/api/routers"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (a *ApiHandler) Handle() {
	// instantiate services

	// create router
	router := chi.NewRouter()

	// setup routes
	// api_routers.SetupComputeRoutes(router, computeService)
	api_routers.SetupHealthRotes(router)

	// start server
	http.ListenAndServe(":8080", router)
}
