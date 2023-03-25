// generated from binary.csgo

package pack

import (
	"apm/csgo/api"
	"github.com/consensys/gnark/frontend"
)

import "github.com/consensys/gnark/std/math/bits"

// ToBinary coverts to...
func ToBinary(bitLen int, x frontend.Variable) []frontend.Variable {
	return bits.ToBinary(api.Api, x, bits.WithNbDigits(bitLen))
}
