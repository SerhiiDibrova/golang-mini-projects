package initializers

import (
	"errors"
	"net/http"
	"strings"
)

type HostVerifier struct {
	allowedHosts []string
}

func NewHostVerifier(allowedHosts []string) *HostVerifier {
	return &HostVerifier{allowedHosts: allowedHosts}
}

func (hv *HostVerifier) Verify(r *http.Request) error {
	host := r.Header.Get("Host")
	if host == "" {
		return errors.New("host header is not present in the request")
	}
	if !hv.allowedHost(host) {
		return errors.New("host is not allowed")
	}
	return nil
}

func (hv *HostVerifier) allowedHost(host string) bool {
	for _, allowedHost := range hv.allowedHosts {
		if strings.ToLower(host) == strings.ToLower(allowedHost) {
			return true
		}
	}
	return false
}