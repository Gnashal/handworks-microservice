package helpers

import (
	"crypto/rand"
	"math/big"
)

type Helpers struct {
}

func (r *Helpers) GenerateIdNum() big.Int {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	return *randNumber
}
