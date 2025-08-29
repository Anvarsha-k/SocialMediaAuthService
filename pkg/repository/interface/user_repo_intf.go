package interfaceRepository_authSvc

// import requestmodels_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/pkg/infrastructure/models/requestmodels"

type IUserRepo interface {
	IsUserExist(email string) bool

	CreateUser(name, email, username, password string) error

	GetHashPassAndStatus(email string) (string, string, string, error)
}
