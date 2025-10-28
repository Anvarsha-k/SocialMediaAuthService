package usecase_authSvc_test

import (
	"fmt"
	"testing"
	"time"

	config_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/config"
	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	usecase_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/// ---------- Mock Definitions ---------- ///

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) IsUserExist(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}
func (m *mockUserRepo) DeleteRecentOtpRequestsBefore5min() error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockUserRepo) TemporarySavingUserOtp(otp int, email string, exp time.Time) error {
	args := m.Called(otp, email, exp)
	return args.Error(0)
}
func (m *mockUserRepo) CreateUser(name, email, username, password string) error {
	args := m.Called(name, email, username, password)
	return args.Error(0)
}

// just mocking to satisfy go
// start
func (m *mockUserRepo) GetHashPassAndStatus(email string) (string, string, string, error) {
	return "", "", "", nil
}
func (m *mockUserRepo) GetOtpInfo(email string) (string, time.Time, error) {
	return "", time.Now(), nil
}
func (m *mockUserRepo) ChangeUserStatusActive(email string) error                  { return nil }
func (m *mockUserRepo) GetUserId(email string) (string, error)                     { return "123", nil }

func (m *mockUserRepo) UpdateUserPassword(email string, hashedPass string) error { 
	args:=m.Called(email,hashedPass)
	return args.Error(0)
}

//end

type mockRandNumUtils struct {
	mock.Mock
}

func (m *mockRandNumUtils) RandomNumber() int {
	args := m.Called()
	return args.Int(0)
}

type mockSendGridUtils struct {
	mock.Mock
}

func (m *mockSendGridUtils) SendVerificationEmailWithOtp(otp int, email string, name string) error {
	args := m.Called(otp, email, name)
	return args.Error(0)
}

type mockJwtUtils struct {
	mock.Mock
}

func (m *mockJwtUtils) TempTokenForOtpVerification(secret, email string) (string, error) {
	args := m.Called(secret, email)
	return args.String(0), args.Error(1)
}

// GenerateAccessToken mocks the method to generate access token
func (m *mockJwtUtils) GenerateAccessToken(secretkey string, id string) (string, error) {
	args := m.Called(secretkey, id)
	return args.String(0), args.Error(1)
}

// GenerateRefreshToken mocks the method to generate refresh token
func (m *mockJwtUtils) GenerateRefreshToken(secretkey string) (string, error) {
	args := m.Called(secretkey)
	return args.String(0), args.Error(1)
}

// UnbindEmailFromClaim mocks the method that extracts email from temp token
func (m *mockJwtUtils) UnbindEmailFromClaim(tokenString string, tempVerificationKey string) (string, error) {
	args := m.Called(tokenString, tempVerificationKey)
	return args.String(0), args.Error(1)
}

//mock these to avoid "not enough argument error" because the constructor of usecase have other arguments

type mockHashUtils struct {
	mock.Mock
}

func (m *mockHashUtils) ComparePassword(hashedpassword, plainPassword string) error {
	args := m.Called(hashedpassword, plainPassword)
	return args.Error(0)
}
func (m *mockHashUtils) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// mock these to avoid "not enough argument error" because the constructor of usecase have other arguments
type mockSmtpUtils struct {
	mock.Mock
}

func (m *mockSmtpUtils) SendResetPasswordEmailOtp(otp int, recieveremail string) error {
	args := m.Called(otp, recieveremail)
	return args.Error(0)
}


/// ---------- Test Starts Here ---------- ///

func TestUserSignUp(t *testing.T) {
	// 1️⃣ Create mocks
	mockRepo := new(mockUserRepo)
	mockRand := new(mockRandNumUtils)
	mockSendGrid := new(mockSendGridUtils)
	mockJwtUtl := new(mockJwtUtils)
	mockHash := new(mockHashUtils)
	mockSmtp := new(mockSmtpUtils)

	// 2️⃣ Setup mock expectations

	mockRepo.On("IsUserExist", "test@example.com").Return(false)
	mockRepo.On("DeleteRecentOtpRequestsBefore5min").Return(nil)
	mockRepo.On("TemporarySavingUserOtp", mock.Anything, "test@example.com", mock.Anything).Return(nil)
	mockRepo.On("CreateUser", "testuser", "test@example.com", "123Test", mock.Anything).Return(nil)

	mockRand.On("RandomNumber").Return(1234)

	mockSendGrid.On("SendVerificationEmailWithOtp", mock.Anything, "test@example.com", "testuser").Return(nil)

	mockJwtUtl.On("TempTokenForOtpVerification", "test_secret", "test@example.com").Return("mocked-temp-token", nil)

	tokenConfig := &config_authSvc.Token{TempVerificationKey: "test_secret"}

	u := usecase_authSvc.NewUserUseCase(mockRepo, mockHash, mockJwtUtl, tokenConfig, mockRand, mockSendGrid, mockSmtp)

	// Now test your function
	req := requestmodels_authSvc.UserSignUpReq{
		Name:            "testuser",
		UserName:        "123Test",
		Email:           "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}

	resp, err := u.UserSignUp(&req)
	fmt.Println(resp.Token)
	fmt.Println(err)//this shows nil because no error signup success and the usersignup return nil and a token
	//🧩 What is assert?
	// assert comes from the testify/assert package.
	// It’s used to verify that certain conditions are true in a test.
	// If an assertion fails — your test fails.

	//only to test the is error situation(error :User already exist with same email)
	//in this situation userexist returns true and the userSignup didnt have revice a resp thatis "_, err := u.UserSignUp(&req)""

	// assert.Error(t,err)
	// assert.EqualError(t,err,"user exist try again with another email")

	//success situation
	assert.NoError(t,err)
	assert.Equal(t, "mocked-temp-token",resp.Token)

}

func TestResetPassword(t *testing.T){
	mockRepo:=new(mockUserRepo)
	mockJwt:=new(mockJwtUtils)
	mockHash:=new(mockHashUtils)
	mockSmtp := new(mockSmtpUtils)
	mockSendGrid := new(mockSendGridUtils)
	mockRand := new(mockRandNumUtils)

	TokenSecurityKey:=config_authSvc.Token{TempVerificationKey: "secret_test"}

	mockJwt.On("UnbindEmailFromClaim",mock.Anything,"secret_test").Return("testing@example.com",nil)
	mockHash.On("HashPassword",mock.Anything).Return("123ABC",nil)
	mockRepo.On("GetOtpInfo","testing@example.com").Return("123456",time.Now(),nil)
	mockRepo.On("UpdateUserPassword","testing@example.com","123ABC").Return(nil)

	u:=usecase_authSvc.NewUserUseCase(mockRepo,mockHash,mockJwt,&TokenSecurityKey,mockRand,mockSendGrid,mockSmtp)

	req:=requestmodels_authSvc.ForgotPasswordData{
		Password: "123ABC",
		ConfirmPassword: "123ABC",
	}

	err:=u.ResetPassword(&req,&TokenSecurityKey.TempVerificationKey)

	assert.NoError(t,err)
}