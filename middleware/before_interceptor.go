package middleware

import (
	"log"
	"net/http"
	"net"
)

func beforeInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		if remoteAddr == "" {
			log.Println("Error: unable to extract remote address")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if remoteAddr == "127.0.0.1" || remoteAddr == "[::1]" {
			log.Printf("INFO: access granted from %s", remoteAddr)
			next.ServeHTTP(w, r)
		} else {
			log.Printf("DEBUG: access denied from %s", remoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}