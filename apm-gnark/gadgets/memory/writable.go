// generated from writable.csgo

package memory

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/selector"
)

type Writable struct {
	keys        []frontend.Variable
	values      []frontend.Variable
	indicators  []frontend.Variable
	writableLen int
	api         frontend.API
}

func NewWritable(api frontend.API, keys, values []frontend.Variable, writableLen int) *Writable {
	return &Writable{api: api, keys: keys, values: values, writableLen: writableLen}
}

func (mem *Writable) AssertValuesAre(values []frontend.Variable) {
	for i := 0; i < len(values); i++ {
		mem.api.AssertIsEqual(mem.values[i], values[i])
	}
}

func (mem *Writable) SelectAddr(addrKey frontend.Variable) {
	mem.indicators = selector.KeyDecoder(mem.api, addrKey, mem.keys)
}

func (mem *Writable) Read(addrKey frontend.Variable) frontend.Variable {
	return selector.Map(mem.api, addrKey, mem.keys, mem.values)
}

func (mem *Writable) Write(wrValue frontend.Variable) {
	for i := 0; i < mem.writableLen; i++ {
		// mem.values[i] <== indicators[i]*(wrValue-mem.values[i]) + mem.values[i]
		mem.values[i] = mem.api.Add(mem.api.Mul(mem.indicators[i], mem.api.Sub(wrValue, mem.values[i])), mem.values[i])
	}
}
