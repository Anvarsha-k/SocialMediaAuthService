package responsemodels_authSvc

type UserSignUpResp struct{
	Token string
}
type UserLoginResp struct{
	AccessToken string
	RefreshToken string
}