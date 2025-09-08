package jwttoken_authSvc

import (
	"fmt"
	"time"

	interface_jwt_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/jwt.go/interface"
	"github.com/golang-jwt/jwt"
)

type JwtUtil struct{}

func NewJwtUtil() interface_jwt_authSvc.IJwt {
	return &JwtUtil{}
}

func (jwtUtil *JwtUtil) GenerateAccessToken(secretkey string, id string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 3600,
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		fmt.Println("Error Creating Access Token!")
		return "", err
	}
	return tokenString, nil

}

func (jwtUtil *JwtUtil) GenerateRefreshToken(secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 604800,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error Creating RefreshToken!")
		return "", err
	}
	return tokenString, nil
}

func (JwtUtil *JwtUtil) TempTokenForOtpVerification(securityKey, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(securityKey))
	if err != nil {
		fmt.Println("Error at Creating Jwt Token!")
	}
	return tokenString,nil
}
