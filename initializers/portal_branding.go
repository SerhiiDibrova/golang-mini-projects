package initializers

import (
	"gorm.io/gorm"
)

type PortalBranding struct {
	ID          uint   `gorm:"primaryKey"`
	FooterContent string `gorm:"type:text"`
}

func GetPortalBranding(db *gorm.DB) (*PortalBranding, error) {
	var portalBranding PortalBranding
	result := db.First(&portalBranding)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &portalBranding, nil
}

func (p *PortalBranding) GetFooterContent() string {
	return p.FooterContent
}