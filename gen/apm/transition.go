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
	var op, memRead, memWrite, jmp, memAddr, instructionOp = decodeInstruction(instruction)

	var read, indicators = s.cache.Read(memAddr)

	var operand, _ = selector.Mux(memRead, instructionOp, read)

	var add = api.Add(s.acc, operand)
	var sub = api.Sub(s.acc, operand)
	var mul = api.Mul(s.acc, operand)
	var less = comparator.IsLess(s.acc, operand)

	var res, _ = selector.Mux(op, add, sub, mul, api.Sub(0, sub))
	s.acc, _ = selector.Mux(jmp, res, s.acc)

	var write, _ = selector.Mux(memWrite, read, res)
	s.cache.Write(write, indicators)

	// both jmp and less are already constrained to be boolean.
	s.pc = api.Add(s.pc, 1)
	s.pc = api.Add(s.pc, api.Mul(api.Mul(jmp, less), api.Sub(instructionOp, s.pc)))
}

func (s *State) AssertOutputIs(values []frontend.Variable) {
	s.cache.AssertValuesAre(values)
}

func decodeInstruction(instruction frontend.Variable) (op, memRead, memWrite, jmp, memAddr, operand frontend.Variable) {
	var parts []frontend.Variable
	parts, _ = api.Compiler().NewHint(decodeHint, 6, instruction)

	pow2_69 := new(big.Int).Lsh(big.NewInt(1), 69)
	api.AssertIsEqual(api.Add(api.Add(api.Add(api.Add(api.Add(parts[0], api.Mul(1<<2, parts[1])), api.Mul(1<<3, parts[2])), api.Mul(1<<4, parts[3])), api.Mul(1<<5, parts[4])), api.Mul(pow2_69, parts[5])), instruction)
	// We need to constrain each part of the decoded instruction to be in the
	// appropriate range. This is done outside of this function:
	// op, memRead, memWrite are constrained by Multiplexers
	// jmp is constrained explicitly
	// memAddr is constrained by the memory module
	// operand does not need to be constrained since it is the last part of the
	// decomposition
	op = parts[0]
	memRead = parts[1]
	memWrite = parts[2]
	jmp = parts[3]
	memAddr = parts[4]
	operand = parts[5]
	return
}

func decodeHint(_ *big.Int, inputs, results []*big.Int) error {
	instruction := inputs[0]
	results[0].And(instruction, big.NewInt(0b11))                                     // op
	results[1].SetUint64(uint64(instruction.Bit(2)))                                  // memRead
	results[2].SetUint64(uint64(instruction.Bit(3)))                                  // memWrite
	results[3].SetUint64(uint64(instruction.Bit(4)))                                  // jmp
	results[4].And(new(big.Int).Rsh(instruction, 5), new(big.Int).SetUint64(1<<64-1)) // memAddr
	results[5].Rsh(instruction, 69)                                                   // operand
	return nil
}

func init() {
	hint.Register(decodeHint)
}
