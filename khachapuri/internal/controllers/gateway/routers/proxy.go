package routers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
)

func SetupProxyRouter(router *chi.Mux, registryService domain.IRegistryService) {
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		parts := strings.Split(host, ".")
		if len(parts) < 2 {
			http.Error(w, "Invalid host", http.StatusBadRequest)
			return
		}
		subdomain := parts[0]

		target, err := registryService.ResolveDomain(subdomain)
		if err != nil {
			http.Error(w, "Domain not found", http.StatusNotFound)
			return
		}

		url, err := url.Parse(target.ServerEndpoint)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	})
}
