package gosmtp_authSvc

import (
	"fmt"
	"net/smtp"

	config_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/config"
	interface_smtp_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/go_smtp/interface"
)

type SmtpCredentials struct {
	SmtpConfig *config_authSvc.Smtp
}

func NewSmtpCredentials(smtpConfig *config_authSvc.Smtp) interface_smtp_authSvc.ISmtp {
	return &SmtpCredentials{SmtpConfig: smtpConfig}
}

func (smtpCredentials *SmtpCredentials) SendResetPasswordEmailOtp(otp int, receiverEmail string) error {

	from := smtpCredentials.SmtpConfig.SmtpSender
	to := []string{receiverEmail}
	password := smtpCredentials.SmtpConfig.SmtpPassword
	smtpHost := smtpCredentials.SmtpConfig.SmtpHost
	smtpPort := smtpCredentials.SmtpConfig.SmtpPort

	subject := "Reset Password - Connexxa"
	body := fmt.Sprintf("Dear %s,\n\nYou recently requested to reset your password for your Connexa account. To complete the process, please use the following One-Time Password (OTP):\n\nOTP: %d\n\nThis OTP is valid for 10 minutes. Please do not share this OTP with anyone for security reasons. If you did not request a password reset, please ignore this email.\n\nThank you,\nThe Ciao Team", receiverEmail, otp)
	message := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

	//Authentication for smtp server
	auth := smtp.PlainAuth("", from, password, smtpHost)

	//Send Message

	err:=smtp.SendMail(smtpHost +":"+ smtpPort,auth,from,to,message)
	if err !=nil{
		fmt.Println("-----",err)
		return err
	}
	return nil
}
