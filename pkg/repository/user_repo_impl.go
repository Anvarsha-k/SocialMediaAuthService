package repository_authSvc

import (

	// requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"

	"errors"
	"fmt"
	"time"

	interfaceRepository_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/repository/interface"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaceRepository_authSvc.IUserRepo {
	return &UserRepository{DB: db}
}

func (d *UserRepository) IsUserExist(email string) bool {
	var userCount int64

	delUnCompletedUser := "DELETE FROM users WHERE email = $1 AND status = $2"
	result := d.DB.Exec(delUnCompletedUser, email, "pending")
	if result.Error != nil {
		fmt.Println("Error in deleting already existing user with the same email and status = pending")
	}

	query := "SELECT COUNT(*) FROM users WHERE email = $1 AND status!= $2"
	err := d.DB.Raw(query, email, "deleted").Row().Scan(&userCount)
	if err != nil {
		fmt.Println("Error in usercount query")
	}

	if userCount >= 1 {
		return true
	}
	return false
}

// Create new user
func (d *UserRepository) CreateUser(name, email, username, password string) error {
	query := `INSERT INTO users (name, email, user_name, password,created_at, updated_at) VALUES ($1,$2,$3,$4, NOW(), NOW())`

	result := d.DB.Exec(query, name, email, username, password)
	return result.Error
}

func (d *UserRepository) GetHashPassAndStatus(email string) (string, string, string, error) {
	var hasedPassword, status, userid string

	query := "SELECT password,status,id FROM users WHERE email = ? AND status !='delete'"
	err := d.DB.Raw(query, email).Row().Scan(&hasedPassword, &status, &userid)
	if err != nil {
		return "", "", "", errors.New("no user exist with this email,signup first")
	}
	return hasedPassword, status, userid, nil
}

func (d *UserRepository) DeleteRecentOtpRequestsBefore5min() error {
	query := "DELETE FROM otp_infos WHERE expiration < CURRENT_TIMESTAMP - INTERVAL '5 minutes';"
	err := d.DB.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepository) TemporarySavingUserOtp(otp int, userEmail string, expiration time.Time) error {
	query := "INSERT INTO otp_infos (email,otp,expiration) VALUES ($1,$2,$3)"
	err := d.DB.Exec(query, userEmail, otp, expiration).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepository) GetOtpInfo(email string) (string, time.Time, error) {
	var expiration time.Time

	type OTPInfo struct {
		Otp        string    `gorm:"column:otp"`
		Expiration time.Time `gorm:"column:expiration"`
	}
	var otpInfo OTPInfo
	if err := d.DB.Raw("SELECT otp, expiration FROM otp_infos WHERE email = ? ORDER BY expiration DESC LIMIT 1;", email).Scan(&otpInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", expiration, errors.New("otp verification failed, no data found on the user email")
		}
		return "", expiration, errors.New("Otp Verification failed, Error finding user data:" + err.Error())
	}
	return otpInfo.Otp, otpInfo.Expiration, nil
}

func (d *UserRepository) ChangeUserStatusActive(email string) error {
	query := "UPDATE users SET status = 'active' WHERE email = $1"
	result := d.DB.Exec(query, email)
	if result.Error != nil {
		fmt.Println("", result.Error)
		return result.Error
	}
	return nil
}

func (d *UserRepository) GetUserId(email string) (string, error) {
	var userId string
	query := "SELECT id FROM users WHERE email=$1 AND status=$2"
	err := d.DB.Raw(query, email, "active").Row().Scan(&userId)
	if err != nil {
		fmt.Println("", err)
		return "", err
	}
	return userId, nil

}

func (d *UserRepository) UpdateUserPassword(email *string, hashedPassword *string) error {
	query := "UPDATE users SET password=$1 WHERE email=$2"
	err := d.DB.Exec(query, *hashedPassword, *email).Error
	if err != nil {
		return err
	}
	return nil
}
