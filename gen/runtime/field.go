package runtime

import (
	"apm/csgo/api"
	"math/big"
)

func FieldOrder() *big.Int {
	return api.Api.Compiler().Field()
}
