package balancer

import (
	"net/url"
	"sync"
)

type LoadBalancer interface {
	Next() *url.URL
}

type RoundRobin struct {
	targets []*url.URL
	index   int 
	mu      sync.Mutex 
}

func NewRoundRobin(targets []*url.URL) *RoundRobin  {
	return &RoundRobin{targets: targets}
}

func (r *RoundRobin) Next() *url.URL {
	r.mu.Lock()
	defer r.mu.Unlock()

	target := r.targets[r.index]
	r.index = (r.index + 1) % len(r.targets)

	return target
}