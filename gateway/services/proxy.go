package services

import (
	"log"
	"net/http/httputil"
	"net/url"
)

type ProxyService struct {
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (*ProxyService) Proxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	return httputil.NewSingleHostReverseProxy(url)
}
