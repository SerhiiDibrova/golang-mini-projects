package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConnector struct {
	db *gorm.DB
}

func NewDatabaseConnector() *DatabaseConnector {
	return &DatabaseConnector{}
}

func (dc *DatabaseConnector) Connect() error {
	if os.Getenv("DB_USERNAME") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_HOST") == "" || os.Getenv("DB_NAME") == "" {
		return nil
	}

	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	maxAttempts := 5
	attempt := 0

	for attempt < maxAttempts {
		dc.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println(err)
			attempt++
			time.Sleep(2 * time.Second)
			continue
		}

		if dc.db.Error == nil {
			break
		}

		log.Println(dc.db.Error)
		attempt++
		time.Sleep(2 * time.Second)
	}

	if attempt == maxAttempts {
		return err
	}

	return nil
}

func (dc *DatabaseConnector) GetDB() *gorm.DB {
	return dc.db
}

func (dc *DatabaseConnector) Close() error {
	sqlDB, err := dc.db.DB()
	if err != nil {
		log.Println(err)
		return err
	}
	return sqlDB.Close()
}