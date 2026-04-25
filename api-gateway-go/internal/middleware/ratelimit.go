package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count int
	last  time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.RWMutex
	limit   = 5
	window  = time.Minute
)

func RateLimit() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			mu.Lock()
			c, exists := clients[ip]
			if !exists {
				c = &client{count: 1, last: time.Now()}
				clients[ip] = c
			} else {
				if time.Since(c.last) > window {
					c.count = 1
					c.last = time.Now()
				} else {
					c.count++
				}
			}

			count := c.count
			mu.Unlock()

			if count > limit {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
