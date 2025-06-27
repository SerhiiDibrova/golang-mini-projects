package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func _addDownloadTokenCookie(w http.ResponseWriter) {
	token, err := generateDownloadToken()
	if err != nil {
		http.Error(w, "Failed to generate download token", http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     "download_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	if err := setCookie(w, cookie); err != nil {
		http.Error(w, "Failed to set download token cookie", http.StatusInternalServerError)
		return
	}
}

func generateDownloadToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(b)
	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

func setCookie(w http.ResponseWriter, cookie *http.Cookie) error {
	if err := cookie.Validate(); err != nil {
		return err
	}
	http.SetCookie(w, cookie)
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_addDownloadTokenCookie(w)
		w.Write([]byte("Download token cookie added"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}