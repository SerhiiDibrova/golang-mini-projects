package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func (d *Database) Connect(dsn string) error {
	var err error
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

func (d *Database) GetFeatureCount() (int64, error) {
	var count int64
	result := d.db.Model(&struct{}{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}