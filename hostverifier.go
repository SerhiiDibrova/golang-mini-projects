package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HostVerifier struct {
	allowedHosts []string
}

func NewHostVerifier(allowedHosts []string) *HostVerifier {
	return &HostVerifier{allowedHosts: allowedHosts}
}

func (hv *HostVerifier) AllowedHost(host string) bool {
	for _, allowedHost := range hv.allowedHosts {
		if host == allowedHost {
			return true
		}
	}
	return false
}

func main() {
	allowedHosts := []string{"example.com", "localhost"}
	hostVerifier := NewHostVerifier(allowedHosts)
	fmt.Println(hostVerifier.AllowedHost("example.com"))  
	fmt.Println(hostVerifier.AllowedHost("localhost"))   
	fmt.Println(hostVerifier.AllowedHost("other.com"))   
}