package interface_jwt_authSvc

type IJwt interface{
	GenerateAccessToken(secretkey string,id string) (string,error)
	GenerateRefreshToken(secretkey string) (string,error)
	TempTokenForOtpVerification(securityKey string,email string) (string, error)
}