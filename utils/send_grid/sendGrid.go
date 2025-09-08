package sendgrid_authSvc

import (
	"fmt"
	config_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/config"
	interface_sendgrid_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/send_grid/interface"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	sendGridConfig *config_authSvc.SendGridConfig
}

func NewSendGrid(sendgridconfig *config_authSvc.SendGridConfig) interface_sendgrid_authSvc.ISendGrid {
	return &SendGridClient{sendGridConfig: sendgridconfig}
}

func (s *SendGridClient) SendVerificationEmailWithOtp(otp int, recieverEmail, recieverName string) error {
	from := mail.NewEmail("Connexa", s.sendGridConfig.SenderEmail)
	to := mail.NewEmail(recieverName, recieverEmail)
	subject := "Verify Your Email Address For Connexa"
	body:= fmt.Sprintf("Hello,%s\n\nThank you for signing up for Connexa. To complete your registration and ensure the security of your account, please verify your email address by entering the One-Time Password (OTP) provided below:\n\nOTP: %d\n\nPlease use the OTP to verify your email address on our platform within the next 10 minutes. After this time, the OTP will expire, and you will need to request a new one.\n\nIf you did not request this verification, please disregard this email.\n\nIf you need any assistance or have questions, feel free to reach out to our support team at support@example.com.\n\nThank you for choosing Ciao.\n\nBest regards,\nThe Connexa Team", recieverName, otp)
	htmlContent := fmt.Sprintf(
        "<p>Hello %s,</p>"+
            "<p>Thank you for signing up for <b>Connexa</b>!</p>"+
            "<p>Your One-Time Password (OTP) is: <strong>%06d</strong></p>"+
            "<p>This OTP will expire in 10 minutes.</p>"+
            "<p>If you didn’t request this, you can safely ignore this email.</p>"+
            "<br><p>Best regards,<br>The Connexa Team</p>",
        recieverName, otp,
    )
	message :=mail.NewSingleEmail(from,subject,to,body,htmlContent)
	client :=sendgrid.NewSendClient(s.sendGridConfig.APIKey)
	_,err:=client.Send(message)
	// log.Printf("SendGrid response: %d, body: %s", resp.StatusCode, resp.Body)
	return err
}

