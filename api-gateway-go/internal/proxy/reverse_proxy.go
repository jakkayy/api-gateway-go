package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"api-gateway-go/internal/balancer"
)

type Proxy struct {
	lb balancer.LoadBalancer
}

func New(lb balancer.LoadBalancer) http.Handler {
	return &Proxy{lb: lb}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	maxRetries := 3
	timeout := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		target := p.lb.Next()

		log.Println("Try:", i+1, "->", target)

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		req := r.Clone(ctx)

		var proxyErr error
		rp := httputil.NewSingleHostReverseProxy(target)
		rp.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
			proxyErr = err
		}

		rp.ServeHTTP(w, req)
		cancel()

		if proxyErr == nil {
			return
		}

		log.Println("Error:", proxyErr)
		time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
	}
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}
