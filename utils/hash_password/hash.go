package hashpassword_authSvc

import (
	"errors"
	"fmt"

	interface_hash_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/hash_password/interface"
	"golang.org/x/crypto/bcrypt"
)

type HashUtil struct{}

func NewHashUtil() interface_hash_authSvc.IhashPassword {
	return &HashUtil{}
}
func (hashUtil *HashUtil) ComparePassword(hashedPassword string, plainPassword string) error {
	err:= bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(plainPassword))

	if err!=nil{
		return errors.New("password doesnot match")
	}
	return nil
}

func (hashUtil *HashUtil)HashPassword(password string)(string,error){
	hashedPass,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err!=nil{
		fmt.Println("Error in hashing password")
	}
	return string(hashedPass),nil
}