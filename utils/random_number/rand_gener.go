package randnumgene_authSvc

import (
	"crypto/rand"
	"log"
	"math/big"

	interface_randnumgene_authSvc "github.com/Anvarsha-k/SocialMediaAuthService/utils/random_number/interface"
)

type RandomNum struct{}

func NewRandomNumUtil() interface_randnumgene_authSvc.IRandGene {
	return &RandomNum{}
}

func (rn RandomNum) RandomNumber() int {
	max := big.NewInt(1000000)
	n,_ :=rand.Int(rand.Reader,max)
	log.Println(n)
	return int(n.Int64())
}
