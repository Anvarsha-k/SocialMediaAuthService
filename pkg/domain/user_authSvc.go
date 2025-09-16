package domain_authSvc

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName  string `gorm:"size:50;not null"`
	Name      string `gorm:"size:100"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Status    string `gorm:"size:20;default:'pending'"` //pending,active
	// CreatedAt time.Time
	// UpdatedAt time.Time
}
