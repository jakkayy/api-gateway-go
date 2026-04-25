package router

import (
	"net/http"
	"net/url"

	"api-gateway-go/internal/balancer"
	"api-gateway-go/internal/config"
	"api-gateway-go/internal/middleware"
	"api-gateway-go/internal/proxy"
	"api-gateway-go/pkg/logger"
)

func NewRouter(cfg config.Config, logg *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	for _, r := range cfg.Routes {
		var urls []*url.URL

		for _, t := range r.Targets {
			u, _ := url.Parse(t)
			urls = append(urls, u)
		}

		lb := balancer.NewRoundRobin(urls)
		p := proxy.New(lb)

		handler := middleware.Chain(
			middleware.Recovery(),
			middleware.Logging(logg),
			middleware.RateLimit(),
		)(p)

		mux.Handle(r.Path, handler)
	}

	return mux
}
