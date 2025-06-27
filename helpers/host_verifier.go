package helpers

import (
	"errors"
	"net/http"
)

type hostVerifier struct {
	allowedHosts []string
}

func newHostVerifier(allowedHosts []string) *hostVerifier {
	return &hostVerifier{allowedHosts: allowedHosts}
}

func (hv *hostVerifier) verifyHost(server *http.Server) (bool, error) {
	if server == nil {
		return false, errors.New("server cannot be nil")
	}
	for _, host := range hv.allowedHosts {
		if server.Addr == host {
			return true, nil
		}
	}
	return false, errors.New("host is not allowed")
}

func allowedHost(server *http.Server, allowedHosts []string) (bool, error) {
	if server == nil {
		return false, errors.New("server cannot be nil")
	}
	if len(allowedHosts) == 0 {
		return false, errors.New("allowed hosts cannot be empty")
	}
	hostVerifier := newHostVerifier(allowedHosts)
	return hostVerifier.verifyHost(server)
}