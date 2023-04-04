// generated from writable.csgo

package memory

import (
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/consensys/gnark/frontend"
)

import "github.com/argennon-project/csgo/transpiled/selector"

type Writable struct {
	keys        []frontend.Variable
	values      []frontend.Variable
	writableLen int
}

func NewWritable(keys, values []frontend.Variable, writableLen int) *Writable {
	return &Writable{keys: keys, values: values, writableLen: writableLen}
}

func (mem *Writable) AssertValuesAre(values []frontend.Variable) {
	for i := 0; i < len(values); i++ {
		api.AssertIsEqual(mem.values[i], values[i])
	}
}

func (mem *Writable) Read(addrKey frontend.Variable) (readValue frontend.Variable, indicators []frontend.Variable) {
	return selector.Map(addrKey, mem.keys, mem.values)
}

func (mem *Writable) Write(wrValue frontend.Variable, indicators []frontend.Variable) {
	if len(indicators) != mem.writableLen {
		panic("invalid indicators")
	}

	for i := 0; i < mem.writableLen; i++ {
		mem.values[i] = api.Add(api.Mul(indicators[i], api.Sub(wrValue, mem.values[i])), mem.values[i])
	}
}
