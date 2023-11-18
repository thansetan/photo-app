package models

import "time"

type Photo struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Caption   string
	PhotoPath string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsPrivate bool

	User User
}
