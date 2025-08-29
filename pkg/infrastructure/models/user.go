package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserName  string `gorm:"size:50;not null"`
	Name      string `gorm:"size:100"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Status    string `gorm:"size:20;default:'pending'"` //pending,active
	CreatedAt time.Time
	UpdatedAt time.Time
}
