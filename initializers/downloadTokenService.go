package initializers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type DownloadTokenService struct{}

func (d *DownloadTokenService) AddDownloadTokenCookie(w http.ResponseWriter, downloadToken string) error {
	if downloadToken == "" {
		return errors.New("download token is empty")
	}

	cookie := &http.Cookie{
		Name:  "download_token",
		Value: downloadToken,
	}

	err := d.setDownloadTokenCookie(w, cookie)
	if err != nil {
		return d.HandleDownloadTokenCookieError(w, err)
	}

	return nil
}

func (d *DownloadTokenService) setDownloadTokenCookie(w http.ResponseWriter, cookie *http.Cookie) error {
	err := d.validateCookie(cookie)
	if err != nil {
		return err
	}

	err = d.addCookieToResponse(w, cookie)
	if err != nil {
		return err
	}

	return nil
}

func (d *DownloadTokenService) validateCookie(cookie *http.Cookie) error {
	if cookie == nil {
		return errors.New("cookie is nil")
	}

	if cookie.Name == "" {
		return errors.New("cookie name is empty")
	}

	if cookie.Value == "" {
		return errors.New("cookie value is empty")
	}

	return nil
}

func (d *DownloadTokenService) addCookieToResponse(w http.ResponseWriter, cookie *http.Cookie) error {
	err := http.SetCookie(w, cookie)
	if err != nil {
		return err
	}

	return nil
}

func (d *DownloadTokenService) HandleDownloadTokenCookieError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}