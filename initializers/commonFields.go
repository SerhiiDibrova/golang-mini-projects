package initializers

import (
	"errors"
	"log"
	"net/url"
	"strings"
)

type CommonFields struct{}

func (cf *CommonFields) _loadCommonFields(inputMap map[string]string) (map[string]string, error) {
	if inputMap == nil {
		return nil, errors.New("input map is nil")
	}

	commonFields := make(map[string]string)

	requiredFields := []string{"field_name", "url_substitution"}

	for _, field := range requiredFields {
		if _, ok := inputMap[field]; !ok {
			return nil, errors.New("input map is missing required field: " + field)
		}
	}

	for key, value := range inputMap {
		if strings.Contains(key, "field_name") {
			commonFields[key] = value
		} else if strings.Contains(key, "url_substitution") {
			parsedUrl, err := url.Parse(value)
			if err != nil {
				log.Printf("error parsing url substitution: %v", err)
				return nil, errors.New("invalid url substitution")
			}
			commonFields[key] = parsedUrl.String()
		}
	}

	return commonFields, nil
}

func NewCommonFields() *CommonFields {
	return &CommonFields{}
}