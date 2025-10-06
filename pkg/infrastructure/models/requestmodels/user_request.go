package requestmodels_authSvc

type UserSignUpReq struct {
	Name            string `json:"userName"`
	UserName        string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
type UserLoginReq struct {
	Email    string
	Password string
}

type 	ForgotPasswordData struct{
	Otp string
	Password string
	ConfirmPassword string
}
