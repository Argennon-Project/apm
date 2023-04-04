// generated from transition.csgo

package apm

import (
	"github.com/argennon-project/csgo/transpiled/gnark/api"
	"github.com/consensys/gnark/backend/hint"
	"github.com/consensys/gnark/frontend"
)

import (
	"apm/gadgets/memory"
	"github.com/argennon-project/csgo/transpiled/cmp"
	"github.com/argennon-project/csgo/transpiled/selector"
	"math/big"
)

var comparator *cmp.BoundedComparator

type State struct {
	instructions []frontend.Variable
	pc           frontend.Variable
	acc          frontend.Variable
	cache        memory.Writable
}

func NewState(instructions []frontend.Variable, cache memory.Writable) *State {
	comparator = cmp.NewBoundedComparator(big.NewInt(1<<32-1), false)
	return &State{
		instructions: instructions,
		pc:           0,
		acc:          0,
		cache:        cache,
	}
}

func (s *State) Transit() {
	var instruction, _ = selector.Mux(s.pc, s.instructions...)
	var op, memRead, memWrite, jmp, instructionSmallOp, instructionBigOp = decodeInstruction(instruction)

	var readAddr, _ = selector.Mux(memRead, 0, instructionBigOp)
	var read, indicators = s.cache.Read(readAddr)

	var operand, _ = selector.Mux(memRead, instructionBigOp, read)

	var add = api.Add(s.acc, operand)
	var sub = api.Sub(s.acc, operand)
	var mul = api.Mul(s.acc, operand)
	var less = comparator.IsLess(s.acc, operand)

	var res, _ = selector.Mux(op, add, sub, mul, api.Sub(0, sub))
	s.acc, _ = selector.Mux(memWrite, res, instructionSmallOp)

	var write, _ = selector.Mux(memWrite, read, res)
	s.cache.Write(write, indicators)

	var jumpAddr, _ = selector.Mux(less, api.Add(s.pc, 1), instructionSmallOp)
	s.pc, _ = selector.Mux(jmp, api.Add(s.pc, 1), jumpAddr)
}

func (s *State) AssertOutputIs(values []frontend.Variable) {
	s.cache.AssertValuesAre(values)
}

func decodeInstruction(instruction frontend.Variable) (op, memRead, memWrite, jmp, smallOperand, bigOperand frontend.Variable) {
	var parts []frontend.Variable
	parts, _ = api.Compiler().NewHint(decodeHint, 6, instruction)

	api.AssertIsEqual(api.Add(api.Add(api.Add(api.Add(api.Add(parts[0], api.Mul(1<<2, parts[1])), api.Mul(1<<3, parts[2])), api.Mul(1<<4, parts[3])), api.Mul(1<<5, parts[4])), api.Mul(1<<21, parts[5])), instruction)

	op = parts[0]
	memRead = parts[1]
	memWrite = parts[2]
	jmp = parts[3]
	smallOperand = parts[4]
	bigOperand = parts[5]
	return
}

func decodeHint(_ *big.Int, inputs, results []*big.Int) error {
	instruction := inputs[0]
	results[0].And(instruction, big.NewInt(0b11))                         // op
	results[1].SetUint64(uint64(instruction.Bit(2)))                      // memRead
	results[2].SetUint64(uint64(instruction.Bit(3)))                      // memWrite
	results[3].SetUint64(uint64(instruction.Bit(4)))                      // jmp
	results[4].And(new(big.Int).Rsh(instruction, 5), big.NewInt(1<<16-1)) // smallOperand
	results[5].Rsh(instruction, 21)                                       // bigOperand
	return nil
}

func init() {
	hint.Register(decodeHint)
}
