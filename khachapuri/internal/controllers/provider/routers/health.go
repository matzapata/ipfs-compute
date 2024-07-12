package routers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func SetupHealthRotes(router *chi.Mux) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
