package interfaceUseCase_authSvc

import (
	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/responsemodels"
)

type IUserUseCase interface {
	UserSignUp(*requestmodels_authSvc.UserSignUpReq) (*responsemodels_authSvc.UserSignUpResp, error)
	UserLogin(*requestmodels_authSvc.UserLoginReq) (responsemodels_authSvc.UserLoginResp, error)
	VerifyOtp(otp string, TempVerificationToken *string) (responsemodels_authSvc.OtpVerifResult, error)
	ForgotPasswordRequest(email *string) (*string, error)
	ResetPassword(userData *requestmodels_authSvc.ForgotPasswordData, TempToken *string)error
}
