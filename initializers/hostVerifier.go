package initializers

import (
	"sync"
)

type HostVerifier struct {
	allowedHosts map[string]bool
	mu           sync.RWMutex
}

func NewHostVerifier() *HostVerifier {
	return &HostVerifier{
		allowedHosts: make(map[string]bool),
	}
}

func (h *HostVerifier) AddAllowedHost(host string) {
	h.mu.Lock()
	h.allowedHosts[host] = true
	h.mu.Unlock()
}

func (h *HostVerifier) RemoveAllowedHost(host string) {
	h.mu.Lock()
	delete(h.allowedHosts, host)
	h.mu.Unlock()
}

func (h *HostVerifier) AllowedHost(host string) bool {
	h.mu.RLock()
	allowed := h.allowedHosts[host]
	h.mu.RUnlock()
	return allowed
}

func (h *HostVerifier) AddAllowedHosts(hosts []string) {
	h.mu.Lock()
	for _, host := range hosts {
		h.allowedHosts[host] = true
	}
	h.mu.Unlock()
}

func (h *HostVerifier) RemoveAllowedHosts(hosts []string) {
	h.mu.Lock()
	for _, host := range hosts {
		delete(h.allowedHosts, host)
	}
	h.mu.Unlock()
}