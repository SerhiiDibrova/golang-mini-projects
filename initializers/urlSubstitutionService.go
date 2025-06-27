package initializers

import (
	"errors"
	"strings"
)

type UrlSubstitutionService struct{}

func (u *UrlSubstitutionService) ApplySubstitutions(url string, urlSubstitutions map[string]string) (string, error) {
	if urlSubstitutions == nil || len(urlSubstitutions) == 0 {
		return "", errors.New("urlSubstitutions map is empty or nil")
	}
	for old, new := range urlSubstitutions {
		url = strings.ReplaceAll(url, old, new)
	}
	return url, nil
}

func NewUrlSubstitutionService() *UrlSubstitutionService {
	service := &UrlSubstitutionService{}
	if service == nil {
		panic("Failed to create UrlSubstitutionService instance")
	}
	return service
}