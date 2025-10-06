package server

import (
	"context"
	"fmt"
	"log"

	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	"github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/pb"
	interfaceUseCase_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/usecase/interface"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userUseCase interfaceUseCase_authSvc.IUserUseCase
}

func NewAuthService(userUserCase interfaceUseCase_authSvc.IUserUseCase) *AuthService {
	return &AuthService{userUseCase: userUserCase}
}

func (s *AuthService) UserSignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	if req.Email == "" || req.Name == "" || req.UserName == "" || req.Password == "" {
		log.Printf("Invalid request: missing required fields")
		return &pb.SignUpResponse{
			ErrorMessage: "All fields are required",
		}, nil
	}

	// Convert gRPC request to internal request model
	inputData := &requestmodels_authSvc.UserSignUpReq{
		Name:            req.Name,
		Email:           req.Email,
		UserName:        req.UserName,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
	log.Printf("Converted to input Data: %+v", inputData)
	result, err := s.userUseCase.UserSignUp(inputData)
	if err != nil {
		return &pb.SignUpResponse{
			ErrorMessage: err.Error(),
		}, err
	}
	return &pb.SignUpResponse{
		Token: result.Token,
	}, nil
}

func (s *AuthService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
    return &pb.PingResponse{Message: "Auth service is alive"}, nil
}

func (s *AuthService) UserLogin(ctx context.Context, req *pb.RequestUserLogin) (*pb.ResponseUserLogin, error) {

	var loginData requestmodels_authSvc.UserLoginReq

	if req.Email == "" && req.Password == "" {
		log.Println("Please input credentials")
		return &pb.ResponseUserLogin{
			ErrorMessage: "All fields are Required",
		}, nil
	}
	loginData.Email = req.Email
	loginData.Password = req.Password

	respData, err := s.userUseCase.UserLogin(&loginData)
	if err != nil {
		return &pb.ResponseUserLogin{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseUserLogin{
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefreshToken,
	}, nil

}

func (s *AuthService) UserOTPVerication(ctx context.Context, req *pb.RequestOtpVefification) (*pb.ResponseOtpVerification, error) {
	respData, err := s.userUseCase.VerifyOtp(req.Otp, &req.TempToken)
	fmt.Println(req.Otp)
	if err != nil {
		return &pb.ResponseOtpVerification{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseOtpVerification{
		Otp:          respData.Otp,
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefreshToken,
	}, nil
}

func (s *AuthService) ForgotPasswordRequest(ctx context.Context, req *pb.RequestForgotPass) (*pb.ResponseForgotPass, error) {
	respData, err := s.userUseCase.ForgotPasswordRequest(&req.Email)
	if err != nil {
		return &pb.ResponseForgotPass{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseForgotPass{
		Token: *respData,
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *pb.RequestResetPass) (*pb.ResponseErrorMessage, error) {

	var requestData requestmodels_authSvc.ForgotPasswordData

	requestData.Otp = req.Otp
	requestData.Password = req.Password
	requestData.ConfirmPassword = req.ConfirmPassword

	err := s.userUseCase.ResetPassword(&requestData, &req.TempToken)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseErrorMessage{}, nil
}

func (s *AuthService)VerifyAccessToken(ctx *context.Context)