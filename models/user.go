package models

import "time"

type User struct {
	ID        string `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Photos    []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
