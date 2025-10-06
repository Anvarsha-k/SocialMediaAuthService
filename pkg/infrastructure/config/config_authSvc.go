package config_authSvc

import (

	"github.com/spf13/viper"
)

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
}
type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBName     string `mapstructure:"DBName"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
}

// type Smtp struct {
// 	SmtpSender string `mapstructure:"SMTP_SENDER"`
// 	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
// 	SmtpHost string `mapstructure:"SMTP_HOST"`
// 	SmtpPort string `mapstructure:"SMTP_PORT"`
// }

type SendGridConfig struct {
	APIKey      string `mapstructure:"SENDGRID_API_KEY"`
	SenderEmail string `mapstructure:"SENDGRID_SENDER"`
}

type Smtp struct{
	SmtpSender string `mapstructure:"SMTP_SENDER"`
	SmtpPassword string `mapstructure:"SMTP_APPKEY"`
	SmtpPort string	`mapstructure:"SMTP_PORT"`
	SmtpHost string	`mapstructure:"SMTP_HOST"`
}

type Token struct {
	UserSecurityKey string `mapstructure:"USER_TOKENKEY"`
	TempVerificationKey string `mapstructure:"TempVery_TOKENKEY"`
}	

type Config struct {
	PortMngr PortManager
	DB       DataBase
	Token    Token
	Smtp 	 Smtp
	SendGrid      SendGridConfig
	EmailProvider string `mapstructure:"EMAIL_PROVIDER"` //SENDGRID
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var db DataBase
	var token Token
	// var smtp Smtp
	var sendgrid SendGridConfig
	var smtp Smtp

	viper.AddConfigPath(".") //curent path
	viper.AddConfigPath("..")//parent path
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&portmngr)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&token)
	if err != nil {
		return nil, err
	}
	// err = viper.Unmarshal(&smtp)
	// if err !=nil{
	// 	return nil,err
	// }
	err = viper.Unmarshal(&sendgrid)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&smtp)
	if err !=nil{
		return nil,err
	}

	config := Config{PortMngr: portmngr, DB: db, Token: token, SendGrid: sendgrid, EmailProvider: viper.GetString("EMAIL_PROVIDER"),Smtp: smtp}
	return &config, nil
}
