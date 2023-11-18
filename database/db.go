package database

import (
	"fmt"
	"photo-app/helpers"
	"photo-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(conf helpers.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Asia/Jakarta",
		conf.Host,
		conf.User,
		conf.Password,
		conf.Name,
		conf.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(models.User{}, models.Photo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
