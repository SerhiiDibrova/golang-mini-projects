package initializers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UrlListForFeatureRequestRepository struct {
	db *gorm.DB
}

func NewUrlListForFeatureRequestRepository(db *gorm.DB) *UrlListForFeatureRequestRepository {
	return &UrlListForFeatureRequestRepository{db: db}
}

func (r *UrlListForFeatureRequestRepository) ProcessPropertyName(propertyName string) (string, error) {
	if propertyName == "" {
		return "", nil
	}
	var result string
	err := r.db.Model(&UrlListForFeatureRequest{}).Where("property_name = ?", propertyName).Find(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *UrlListForFeatureRequestRepository) ProcessUrlList(urlList []string) ([]string, error) {
	if len(urlList) == 0 {
		return nil, nil
	}
	var results []string
	for _, url := range urlList {
		var result string
		err := r.db.Model(&UrlListForFeatureRequest{}).Where("url = ?", url).Find(&result).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *UrlListForFeatureRequestRepository) ProxyIfAllowed(url string, allowProxy bool) (string, error) {
	if !allowProxy {
		return "", nil
	}
	var result string
	err := r.db.Model(&UrlListForFeatureRequest{}).Where("url = ?", url).Find(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *UrlListForFeatureRequestRepository) HandleError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	log.Println(err)
}

type UrlListForFeatureRequest struct {
	PropertyName string `gorm:"column:property_name"`
	Url          string `gorm:"column:url"`
}