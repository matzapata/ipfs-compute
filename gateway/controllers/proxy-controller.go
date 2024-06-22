package controllers

import (
	"net/http"
	"strings"

	"github.com/matzapata/ipfs-compute/gateway/services"
	"github.com/matzapata/ipfs-compute/shared/registry"
)

type ProxyController struct {
	RegistryService *registry.RegistryService
	ProxyService    *services.ProxyService
}

func NewProxyController(registryService *registry.RegistryService, proxyService *services.ProxyService) *ProxyController {
	return &ProxyController{
		RegistryService: registryService,
		ProxyService:    proxyService,
	}
}

func (c *ProxyController) Proxy(w http.ResponseWriter, r *http.Request) {
	// Extract subdomain
	host := r.Host
	parts := strings.Split(host, ".")
	if len(parts) < 3 {
		http.Error(w, "Invalid host", http.StatusBadRequest)
		return
	}
	subdomain := parts[0]

	// registry resolve domain
	resolver, err := c.RegistryService.ResolveDomain(subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	target, err := resolver.Server(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxy := c.ProxyService.Proxy(target)
	proxy.ServeHTTP(w, r)
}
