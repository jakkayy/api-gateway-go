package proxy

import (
	"net/http"
	"net/http/httputil"

	"api-gateway-go/internal/balancer"
)

type Proxy struct {
	lb balancer.LoadBalancer
}

func New(lb balancer.LoadBalancer) http.Handler {
	return &Proxy{lb: lb}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := p.lb.Next()

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
