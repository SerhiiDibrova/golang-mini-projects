package initializers

import (
	"errors"
	"fmt"
	"net/http"
)

type SearchService struct{}

func (s *SearchService) ConstructRedirectURL(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("uuid is required")
	}
	baseURL := "/search/"
	redirectURL := fmt.Sprintf("%s%s", baseURL, uuid)
	return redirectURL, nil
}

func (s *SearchService) GetRedirectURL(w http.ResponseWriter, r *http.Request, uuid string) {
	redirectURL, err := s.ConstructRedirectURL(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}