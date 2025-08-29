package usecase_authSvc

import (
	"errors"
	"log"

	requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/responsemodels"
	interfaceRepository_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/repository/interface"
	interface_hash_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/hash_password/interface"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// userUseCase implements ucinterface.IUserUseCase
type userUseCase struct {
	userRepo  interfaceRepository_authSvc.IUserRepo // repository injected via constructor
	hashUtils interface_hash_authSvc.IhashPassword
}

func NewUserUseCase(userRepo interfaceRepository_authSvc.IUserRepo) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) UserSignUp(rq *requestmodels_authSvc.UserSignUpReq) (*responsemodels_authSvc.UserSignUpResp, error) {

	var resSignup responsemodels_authSvc.UserSignUpResp

	// store user in repository (pass the hashed password)
	if isUserExist := u.userRepo.IsUserExist(rq.Email); isUserExist {
		return &resSignup, errors.New("user exist try again with another email")
	}
	// hash the password with bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(rq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	log.Printf("Attempting to create user with email: %s", rq.Email)

	err = u.userRepo.CreateUser(rq.Name, rq.Email, rq.UserName, string(hashed))
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	log.Printf("User doesnt exist, proceeding with user creation")

	token := uuid.NewString()

	return &responsemodels_authSvc.UserSignUpResp{Token: token}, nil

}
func (u *userUseCase) UserLogin(rq *requestmodels_authSvc.UserLoginReq) (*responsemodels_authSvc.UserLoginResp, error) {
	var resLogin responsemodels_authSvc.UserLoginResp

	hashedpass, status, userid, err := u.userRepo.GetHashPassAndStatus(rq.Email)
	if err != nil {
		return &resLogin, err
	}

	passwordErr := u.hashUtils.ComparePassword(hashedpass, rq.Password)
	if passwordErr != nil {
		return &resLogin, passwordErr
	}

	if status == "blocked" {
		return &resLogin, errors.New("user is blocked by Admin")
	}
	if status == "pending" {
		return &resLogin, errors.New("user is in pending,OTP not verified")
	}
	
}
