package interfaceRepository_authSvc

import "time"

// import requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"

type IUserRepo interface {
	IsUserExist(email string) bool

	CreateUser(name, email, username, password string) error

	GetHashPassAndStatus(email string) (string, string, string, error)

	DeleteRecentOtpRequestsBefore5min() error

	TemporarySavingUserOtp(otp int, email string, expiration time.Time) error

	GetOtpInfo(email string) (string, time.Time, error)

	ChangeUserStatusActive(email string) error

	GetUserId(email string) (string, error)

	UpdateUserPassword(email string, hashedPass string)error
}
