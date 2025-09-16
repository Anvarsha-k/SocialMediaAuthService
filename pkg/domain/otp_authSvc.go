package domain_authSvc

import (
	"time"

	"gorm.io/gorm"
)

type OtpInfo struct {
	gorm.Model
	Email      string
	OTP        int
	Expiration time.Time
}
