package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	ID       uint   `gorm:"primarykey"`
	Setting  string `gorm:"type:varchar(255)"`
	Value    string `gorm:"type:varchar(255)"`
}

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(dsn string) *ConfigRepository {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) GetConfig() ([]Config, error) {
	var configs []Config
	result := r.db.Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

func (r *ConfigRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (r *ConfigRepository) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *ConfigRepository) ValidateConnection() error {
	err := r.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (r *ConfigRepository) GetConfigWithValidation() ([]Config, error) {
	err := r.ValidateConnection()
	if err != nil {
		return nil, err
	}
	return r.GetConfig()
}