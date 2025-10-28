package usecase_authSvc

import (
	"errors"
	"fmt"
	"log"
	"time"

	config_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/config"
	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/responsemodels"
	interfaceRepository_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/repository/interface"
	interface_smtp_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/go_smtp/interface"
	interface_hash_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/hash_password/interface"
	interface_jwt_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/jwt.go/interface"
	interface_randnumgene_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/random_number/interface"
	interface_sendgrid_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/send_grid/interface"
	"golang.org/x/crypto/bcrypt"
)

// userUseCase implements ucinterface.IUserUseCase
type userUseCase struct {
	userRepo         interfaceRepository_authSvc.IUserRepo // repository injected via constructor
	hashUtils        interface_hash_authSvc.IhashPassword
	jwtUtil          interface_jwt_authSvc.IJwt
	TokenSecurityKey *config_authSvc.Token
	RandNumUtils     interface_randnumgene_authSvc.IRandGene
	sendGridUtils    interface_sendgrid_authSvc.ISendGrid
	smtpUtils        interface_smtp_authSvc.ISmtp
}

func NewUserUseCase(userRepo interfaceRepository_authSvc.IUserRepo, hashUtils interface_hash_authSvc.IhashPassword, jwtUtils interface_jwt_authSvc.IJwt, tokenConfig *config_authSvc.Token, randNumUtils interface_randnumgene_authSvc.IRandGene, sendGridUtil interface_sendgrid_authSvc.ISendGrid,smtpUtils interface_smtp_authSvc.ISmtp) *userUseCase {

	return &userUseCase{
		userRepo:         userRepo,
		hashUtils:        hashUtils,
		jwtUtil:          jwtUtils,
		TokenSecurityKey: tokenConfig,
		RandNumUtils:     randNumUtils,
		sendGridUtils:    sendGridUtil,
		smtpUtils:        smtpUtils,
	}
}

func (u *userUseCase) UserSignUp(rq *requestmodels_authSvc.UserSignUpReq) (*responsemodels_authSvc.UserSignUpResp, error) {

	var resSignup responsemodels_authSvc.UserSignUpResp

	// store user in repository (pass the hashed password)
	if isUserExist := u.userRepo.IsUserExist(rq.Email); isUserExist {
		return &resSignup, errors.New("user exist try again with another email")
	}

	//Delete Recent OTP before 5min
	errRem := u.userRepo.DeleteRecentOtpRequestsBefore5min()
	if errRem != nil {
		return &resSignup, errRem
	}
	//otp Generating
	otp := u.RandNumUtils.RandomNumber()
	errOtp := u.sendGridUtils.SendVerificationEmailWithOtp(otp, rq.Email, rq.Name) //sending otp through email to user
	if errOtp != nil {
		log.Printf("Error Sending OTP to eamil: %s", errOtp)
		return &resSignup, errOtp
	}

	expiration := time.Now().Add(5 * time.Minute)

	//Temporary OTP Saving for OTP verification
	errTempSave := u.userRepo.TemporarySavingUserOtp(otp, rq.Email, expiration)
	if errTempSave != nil {
		fmt.Println("Cant save temporary data for otp verification in db")
		return &resSignup, errors.New("OTP verification down,please try after some time")
	}

	// hash the password with bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(rq.Password), bcrypt.DefaultCost)
	if err != nil {
		return &resSignup, err
	}
	log.Printf("Attempting to create user with email: %s", rq.Email)

	err = u.userRepo.CreateUser(rq.Name, rq.Email, rq.UserName, string(hashed))
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return &resSignup, err
	}
	log.Printf("User doesnt exist, proceeding with user creation")

	tempToken, err := u.jwtUtil.TempTokenForOtpVerification(u.TokenSecurityKey.TempVerificationKey, rq.Email)
	if err != nil {
		fmt.Println("error creating temp token for otp verification")
		return &resSignup, errors.New("error creating temp token for otp verification")
	}
	resSignup.Token = tempToken
	return &resSignup, nil

	// token := uuid.NewString()

	// return &responsemodels_authSvc.UserSignUpResp{Token: token}, nil

}

func (u *userUseCase) UserLogin(rq *requestmodels_authSvc.UserLoginReq) (responsemodels_authSvc.UserLoginResp, error) {
	var resLogin responsemodels_authSvc.UserLoginResp

	hashedpass, status, userid, err := u.userRepo.GetHashPassAndStatus(rq.Email)
	if err != nil {
		return resLogin, err
	}

	passwordErr := u.hashUtils.ComparePassword(hashedpass, rq.Password)
	if passwordErr != nil {
		return resLogin, passwordErr
	}

	if status == "blocked" {
		return resLogin, errors.New("user is blocked by Admin")
	}
	if status == "pending" {
		return resLogin, errors.New("user is in pending,OTP not verified")
	}
	//Creating AccessToken
	accessToken, err := u.jwtUtil.GenerateAccessToken(u.TokenSecurityKey.UserSecurityKey, userid)
	if err != nil {
		return resLogin, err
	}

	//Creating RefreshToken

	refreshToken, err := u.jwtUtil.GenerateRefreshToken(u.TokenSecurityKey.UserSecurityKey)
	if err != nil {
		return resLogin, err
	}
	resLogin.AccessToken = accessToken
	resLogin.RefreshToken = refreshToken

	return resLogin, nil
}

func (u *userUseCase) VerifyOtp(otp string, TempVerificationToken *string) (responsemodels_authSvc.OtpVerifResult, error) {
	var otpVeriRes responsemodels_authSvc.OtpVerifResult
	email, unbindErr := u.jwtUtil.UnbindEmailFromClaim(*TempVerificationToken, u.TokenSecurityKey.TempVerificationKey)
	if unbindErr != nil {
		return otpVeriRes, unbindErr
	}
	userOtp, otpExpire, errGetInfo := u.userRepo.GetOtpInfo(email)
	if errGetInfo != nil {
		return otpVeriRes, errGetInfo
	}
	if otp != userOtp {
		return otpVeriRes, errors.New("invalid OTP")
	}
	if time.Now().After(otpExpire) {
		return otpVeriRes, errors.New("OTP Expired")
	}
	changeStatErr := u.userRepo.ChangeUserStatusActive(email)
	if changeStatErr != nil {
		return otpVeriRes, changeStatErr
	}
	UserId, fetchErr := u.userRepo.GetUserId(email)
	if fetchErr != nil {
		return otpVeriRes, fetchErr
	}
	accessToken, aTokenErr := u.jwtUtil.GenerateAccessToken(u.TokenSecurityKey.UserSecurityKey, UserId)
	if aTokenErr != nil {
		return otpVeriRes, aTokenErr
	}
	refreshToken, rToekenErr := u.jwtUtil.GenerateRefreshToken(u.TokenSecurityKey.UserSecurityKey)
	if rToekenErr != nil {
		return otpVeriRes, rToekenErr
	}
	otpVeriRes.AccessToken = accessToken
	otpVeriRes.RefreshToken = refreshToken
	otpVeriRes.Otp = "verified"

	return otpVeriRes, nil

}
func (u *userUseCase) ForgotPasswordRequest(email *string) (*string, error) {
	_, _, status, err := u.userRepo.GetHashPassAndStatus(*email)
	if err != nil {
		return nil, err
	}
	if status == "pending" {
		return nil, errors.New("user status is on pending,OTP not verified")
	}
	if status == "blocked" {
		return nil, errors.New("User is blocked by the Admin")
	}
	err = u.userRepo.DeleteRecentOtpRequestsBefore5min()
	if err != nil {
		return nil, err
	}
	otp := u.RandNumUtils.RandomNumber()
	err = u.smtpUtils.SendResetPasswordEmailOtp(otp, *email)
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(5 * time.Minute)

	errTempSave := u.userRepo.TemporarySavingUserOtp(otp, *email, expiration)
	if errTempSave != nil {
		fmt.Println("Cant save temporary data for otp verification in db")
		return nil, errors.New("OTP verification down,please try after some time")
	}
	tempToken, err := u.jwtUtil.TempTokenForOtpVerification(u.TokenSecurityKey.TempVerificationKey, *email)
	if err != nil {
		return nil, err
	}
	return &tempToken, nil

}

func (u *userUseCase) ResetPassword(userData *requestmodels_authSvc.ForgotPasswordData, TempVerificationToken *string) error {
	email, err := u.jwtUtil.UnbindEmailFromClaim(*TempVerificationToken, u.TokenSecurityKey.TempVerificationKey)
	if err != nil {
		return err
	}
	userOtp, expiration, err := u.userRepo.GetOtpInfo(email)
	if err != nil {
		return err
	}
	if userData.Otp != userOtp {
		return errors.New("invalid OTP!")
	}
	if time.Now().After(expiration) {
		return errors.New("oTP Expired!")
	}
	hashedpass, err := u.hashUtils.HashPassword(userData.Password)
	if err!=nil{
		return err
	}
	err =u.userRepo.UpdateUserPassword(email,hashedpass)
	if err!=nil{
		return err
	}
	return nil
}
