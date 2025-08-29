package interface_hash_authSvc

type IhashPassword interface{
	ComparePassword(hashedPassword string,plainPassword string) error
}