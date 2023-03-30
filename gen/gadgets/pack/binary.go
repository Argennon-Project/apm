// generated from binary.csgo

package pack

import (
	"apm/csgo/api"
	"github.com/consensys/gnark/frontend"
)

import "github.com/consensys/gnark/std/math/bits"

// AssertBitLen ensures that the unsigned binary representation of x has less than bitLen bits.
func AssertBitLen(bitLen int, x frontend.Variable) {
	bits.ToBinary(api.Api, x, bits.WithNbDigits(bitLen))
}
