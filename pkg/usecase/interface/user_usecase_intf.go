package interfaceUseCase_authSvc

import (
	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/responsemodels"
)

type IUserUseCase interface {
	UserSignUp(*requestmodels_authSvc.UserSignUpReq) (*responsemodels_authSvc.UserSignUpResp, error)
	UserLogin(*requestmodels_authSvc.UserLoginReq) (responsemodels_authSvc.UserLoginResp,error)
}
