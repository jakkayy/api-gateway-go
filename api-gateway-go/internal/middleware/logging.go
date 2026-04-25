package middleware

import (
	"net/http"
	"time"

	"api-gateway-go/pkg/logger"
)

func Logging(logg *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logg.Info(r.Method, r.URL.Path, time.Since(start))
		})
	}
}
