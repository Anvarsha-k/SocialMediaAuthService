package interface_smtp_authSvc

type ISmtp interface{
	SendResetPasswordEmailOtp(otp int,receiverEmail string)error
}