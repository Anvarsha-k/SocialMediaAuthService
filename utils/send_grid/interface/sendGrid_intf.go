package interface_sendgrid_authSvc

type ISendGrid interface{
	SendVerificationEmailWithOtp(otp int, recieverEmail, recieverName string)error
}